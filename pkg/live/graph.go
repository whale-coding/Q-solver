package live

import (
	"Q-Solver/pkg/config"
	"Q-Solver/pkg/llm"
	"Q-Solver/pkg/logger"
	"Q-Solver/pkg/prompts"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

// GraphNode 导图节点（发送给前端）
type GraphNode struct {
	ID       string `json:"id"`
	PID      string `json:"pid,omitempty"` // 父节点ID
	Title    string `json:"title"`         // 简短标题
	Question string `json:"question"`      // 完整问题
	Answer   string `json:"answer"`        // 完整回答
}

// Graph 问题导图处理器
type Graph struct {
	ctx           context.Context
	configManager *config.ConfigManager
	llmService    *llm.Service
	emitEvent     func(string, ...any)

	// 消息缓冲
	messages []llm.Message
	mu       sync.Mutex

	// 配置
	triggerRound int // 每多少轮触发一次总结
	currentRound int // 当前轮数

	// 已生成的节点
	nodes []GraphNode
}

// NewGraph 创建问题导图处理器
// triggerRound: 每多少轮对话触发一次总结
func NewGraph(
	ctx context.Context,
	configManager *config.ConfigManager,
	llmService *llm.Service,
	emitEvent func(string, ...any),
	triggerRound int,
) *Graph {
	if triggerRound <= 0 {
		triggerRound = 3
	}
	return &Graph{
		ctx:           ctx,
		configManager: configManager,
		llmService:    llmService,
		emitEvent:     emitEvent,
		messages:      make([]llm.Message, 0),
		nodes:         make([]GraphNode, 0),
		triggerRound:  triggerRound,
	}
}

// Push 推送一轮对话（问题+回答）
func (g *Graph) Push(question, answer string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// 添加到消息缓冲
	g.messages = append(g.messages,
		llm.NewUserMessage(question),
		llm.NewAssistantMessage(answer),
	)
	g.currentRound++

	logger.Printf("Graph: 收到第 %d 轮对话", g.currentRound)

	// 检查是否达到触发轮数
	if g.currentRound >= g.triggerRound {
		// 复制消息用于总结
		msgs := make([]llm.Message, len(g.messages))
		copy(msgs, g.messages)

		// 清空缓冲
		g.messages = make([]llm.Message, 0)
		g.currentRound = 0

		// 异步触发总结
		go g.summarize(msgs)
	}
}

// Clear 清空导图
func (g *Graph) Clear() {
	g.mu.Lock()
	g.messages = make([]llm.Message, 0)
	g.nodes = make([]GraphNode, 0)
	g.currentRound = 0
	g.mu.Unlock()

	g.emitEvent("graph:clear", nil)
}

// summarize 调用辅助模型总结对话
func (g *Graph) summarize(messages []llm.Message) {
	cfg := g.configManager.Get()

	// 检查是否配置了辅助模型
	if cfg.AssistantModel == "" {
		logger.Println("Graph: 未配置辅助模型，跳过总结")
		return
	}

	logger.Printf("Graph: 开始总结 %d 条消息", len(messages))

	// 构建 prompt
	prompt := g.buildPrompt(messages)

	// 调用模型
	ctx, cancel := context.WithTimeout(g.ctx, 30*time.Second)
	defer cancel()

	provider := g.llmService.GetProvider()
	response, err := provider.GenerateContent(ctx, cfg.AssistantModel, []llm.Message{
		llm.NewUserMessage(prompt),
	})
	if err != nil {
		logger.Printf("Graph: 总结失败: %v", err)
		return
	}

	// 解析并添加节点
	nodes := g.parseResponse(response.Content, messages)
	for _, node := range nodes {
		g.mu.Lock()
		g.nodes = append(g.nodes, node)
		g.mu.Unlock()

		g.emitEvent("graph:add-node", node)
		logger.Printf("Graph: 添加节点: %s", node.Title)
	}
}

// buildPrompt 构建提示词
func (g *Graph) buildPrompt(messages []llm.Message) string {
	var sb strings.Builder

	for i := 0; i < len(messages); i += 2 {
		if i+1 < len(messages) {
			sb.WriteString(fmt.Sprintf("\n问题：%s\n回答：%s\n", messages[i].Content, messages[i+1].Content))
		}
	}

	return fmt.Sprintf(prompts.GraphSummarizePromptTemplate, sb.String())
}

// parseResponse 解析模型响应
func (g *Graph) parseResponse(response string, messages []llm.Message) []GraphNode {
	response = strings.TrimSpace(response)

	// 移除 markdown 代码块
	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
		response = strings.TrimSuffix(response, "```")
		response = strings.TrimSpace(response)
	} else if strings.HasPrefix(response, "```") {
		response = strings.TrimPrefix(response, "```")
		response = strings.TrimSuffix(response, "```")
		response = strings.TrimSpace(response)
	}

	var results []struct {
		Title    string `json:"title"`
		Question string `json:"question"`
		Answer   string `json:"answer"`
	}

	if err := json.Unmarshal([]byte(response), &results); err != nil {
		logger.Printf("Graph: 解析JSON失败: %v", err)
		// 解析失败，使用简单模式
		return g.createSimpleNodes(messages)
	}

	// 获取最后一个节点ID作为父节点
	g.mu.Lock()
	var lastNodeID string
	if len(g.nodes) > 0 {
		lastNodeID = g.nodes[len(g.nodes)-1].ID
	}
	g.mu.Unlock()

	nodes := make([]GraphNode, 0, len(results))
	for _, r := range results {
		nodeID := fmt.Sprintf("node-%d", time.Now().UnixNano())
		nodes = append(nodes, GraphNode{
			ID:       nodeID,
			PID:      lastNodeID,
			Title:    r.Title,
			Question: r.Question,
			Answer:   r.Answer,
		})
		lastNodeID = nodeID // 后续节点挂到这个节点下
	}

	return nodes
}

// createSimpleNodes 简单模式生成节点
func (g *Graph) createSimpleNodes(messages []llm.Message) []GraphNode {
	g.mu.Lock()
	var lastNodeID string
	if len(g.nodes) > 0 {
		lastNodeID = g.nodes[len(g.nodes)-1].ID
	}
	g.mu.Unlock()

	nodes := make([]GraphNode, 0)
	for i := 0; i < len(messages); i += 2 {
		if i+1 >= len(messages) {
			break
		}

		question := messages[i].Content
		answer := messages[i+1].Content

		// 截取标题
		title := question
		runes := []rune(title)
		if len(runes) > 15 {
			title = string(runes[:15]) + "..."
		}

		nodeID := fmt.Sprintf("node-%d", time.Now().UnixNano())
		nodes = append(nodes, GraphNode{
			ID:       nodeID,
			PID:      lastNodeID,
			Title:    title,
			Question: question,
			Answer:   answer,
		})
		lastNodeID = nodeID
	}

	return nodes
}
