package llm

import (
	"Q-Solver/pkg/config"
	"Q-Solver/pkg/logger"
	"context"
	"fmt"
	"strings"
)

// CustomAdapter 自定义适配器（作为路由使用）
type CustomAdapter struct {
	config *config.Config
	openai *OpenAIAdapter
	claude *ClaudeAdapter
	gemini *GeminiAdapter
}

// NewCustomAdapter 创建 Custom 适配器
func NewCustomAdapter(cfg *config.Config) *CustomAdapter {
	adapter := &CustomAdapter{
		config: cfg,
		openai: NewOpenAIAdapter(cfg),
		claude: NewClaudeAdapter(cfg),
	}

	// 尝试初始化 Gemini，即使失败也不影响其他适配器
	gemini, err := NewGeminiAdapter(cfg)
	if err != nil {
		logger.Println("CustomAdapter: 初始化 Gemini 适配器失败:", err)
	} else {
		adapter.gemini = gemini
	}

	return adapter
}

// GenerateContentStream 生成内容流
func (a *CustomAdapter) GenerateContentStream(ctx context.Context, messages []Message, callback StreamCallback) (Message, error) {
	model := strings.ToLower(a.config.Model)

	switch {
	case strings.HasPrefix(model, "gemini"):
		if a.gemini == nil {
			return Message{}, fmt.Errorf("Gemini 适配器未初始化")
		}
		logger.Println("CustomAdapter: 路由到 Gemini ->", model)
		return a.gemini.GenerateContentStream(ctx, messages, callback)

	case strings.HasPrefix(model, "claude"):
		logger.Println("CustomAdapter: 路由到 Claude ->", model)
		return a.claude.GenerateContentStream(ctx, messages, callback)

	default:
		logger.Println("CustomAdapter: 路由到 OpenAI ->", model)
		return a.openai.GenerateContentStream(ctx, messages, callback)
	}
}

// TestChat 测试连通性
func (a *CustomAdapter) TestChat(ctx context.Context) error {
	model := strings.ToLower(a.config.Model)

	switch {
	case strings.HasPrefix(model, "gemini"):
		if a.gemini == nil {
			return fmt.Errorf("Gemini 适配器未初始化")
		}
		logger.Println("CustomAdapter: 路由到 Gemini(测试连通性) ->", model)
		return a.gemini.TestChat(ctx)

	case strings.HasPrefix(model, "claude"):
		logger.Println("CustomAdapter: 路由到 Claude(测试连通性) ->", model)
		return a.claude.TestChat(ctx)

	default:
		logger.Println("CustomAdapter: 路由到 OpenAI(测试连通性) ->", model)
		return a.openai.TestChat(ctx)
	}
}

// GenerateContent 非流式生成内容
func (a *CustomAdapter) GenerateContent(ctx context.Context, model string, messages []Message) (Message, error) {
	// 如果指定了模型，使用指定的模型；否则使用配置中的模型
	routeModel := model
	if routeModel == "" {
		routeModel = a.config.Model
	}
	routeModelLower := strings.ToLower(routeModel)

	switch {
	case strings.HasPrefix(routeModelLower, "gemini"):
		if a.gemini == nil {
			return Message{}, fmt.Errorf("Gemini 适配器未初始化")
		}
		logger.Println("CustomAdapter: 路由到 Gemini(非流式) ->", routeModel)
		return a.gemini.GenerateContent(ctx, model, messages)

	case strings.HasPrefix(routeModelLower, "claude"):
		logger.Println("CustomAdapter: 路由到 Claude(非流式) ->", routeModel)
		return a.claude.GenerateContent(ctx, model, messages)

	default:
		logger.Println("CustomAdapter: 路由到 OpenAI(非流式) ->", routeModel)
		return a.openai.GenerateContent(ctx, model, messages)
	}
}

// GetModels 获取模型列表
func (a *CustomAdapter) GetModels(ctx context.Context) ([]string, error) {
	// Custom 模式下，通常连接的是 OneAPI 等聚合层，它们通常兼容 OpenAI 协议
	// 所以默认使用 OpenAI 的 GetModels
	return a.openai.GetModels(ctx)
}
