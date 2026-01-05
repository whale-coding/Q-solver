package solution

import (
	"Q-Solver/pkg/config"
	"Q-Solver/pkg/llm"
	"Q-Solver/pkg/logger"
	"bytes"
	"context"
	"errors"
)

type Callbacks struct {
	EmitEvent func(event string, data ...interface{})
}

type Request struct {
	Config           config.Config
	ScreenshotBase64 string
	ResumeBase64     string
}

type Solver struct {
	llmProvider llm.Provider
	chatHistory []llm.Message // 改用统一的 Message 类型
}

func NewSolver(provider llm.Provider) *Solver {
	return &Solver{
		llmProvider: provider,
		chatHistory: make([]llm.Message, 0),
	}
}

func (s *Solver) SetProvider(provider llm.Provider) {
	s.llmProvider = provider
}

func (s *Solver) ClearHistory() {
	s.chatHistory = make([]llm.Message, 0)
}

func (s *Solver) Solve(ctx context.Context, req Request, cb Callbacks) bool {
	// 1. 检查 API Key
	if req.Config.GetAPIKey() == "" {
		if cb.EmitEvent != nil {
			cb.EmitEvent("require-login")
		}
		return false
	}

	logger.Println("开始解题流程...")

	// 2. 构建 System Prompt
	var systemPrompt bytes.Buffer
	if req.Config.GetPrompt() != "" {
		systemPrompt.WriteString(req.Config.GetPrompt())
	}

	// 如果使用 Markdown 简历，将简历内容追加到 System Prompt
	if req.Config.GetUseMarkdownResume() && req.Config.GetResumeContent() != "" {
		logger.Println("使用 Markdown 简历内容")
		systemPrompt.WriteString("\n\n# 候选人简历内容如下: \n")
		systemPrompt.WriteString(req.Config.GetResumeContent())
	}

	// 3. 构建当前用户消息（包含截图）
	userParts := []llm.ContentPart{
		llm.ImagePart(req.ScreenshotBase64),
	}

	// 如果使用 PDF 简历，将简历附件加入用户消息
	if !req.Config.GetUseMarkdownResume() && req.ResumeBase64 != "" {
		userParts = append(userParts,
			llm.TextPart("\n\n# 候选人简历已作为附件发送，请参考简历内容回答。"),
			llm.PDFPart(req.ResumeBase64),
		)
		logger.Println("已注入简历附件 (PDF)")
	}

	currentUserMsg := llm.NewMultiPartMessage(llm.RoleUser, userParts)

	// 4. 构建最终发送的消息列表
	var messagesToSend []llm.Message

	if req.Config.GetKeepContext() {
		// 保持上下文模式：使用并更新历史记录
		s.ensureSystemPrompt(systemPrompt.String())
		messagesToSend = append(messagesToSend, s.chatHistory...)
	} else {
		// 不保持上下文模式：每次都是全新对话
		messagesToSend = append(messagesToSend, llm.NewSystemMessage(systemPrompt.String()))
	}
	messagesToSend = append(messagesToSend, currentUserMsg)

	// 5. 调用 LLM 生成回答
	if cb.EmitEvent != nil {
		cb.EmitEvent("solution-stream-start")
	}

	response, err := s.llmProvider.GenerateContentStream(ctx, messagesToSend, func(chunk llm.StreamChunk) {
		if cb.EmitEvent != nil {
			// 根据 chunk 类型发送不同事件
			switch chunk.Type {
			case llm.ChunkThinking:
				cb.EmitEvent("solution-stream-thinking", chunk.Content)
			case llm.ChunkContent:
				cb.EmitEvent("solution-stream-chunk", chunk.Content)
			}
		}
	})

	if err != nil {
		if errors.Is(ctx.Err(), context.Canceled) {
			logger.Println("当前任务已中断 (用户产生新输入)")
			if cb.EmitEvent != nil {
				cb.EmitEvent("solution-error", "context canceled")
			}
			return false
		}
		logger.Printf("LLM 请求失败: %v\n", err)
		if cb.EmitEvent != nil {
			cb.EmitEvent("solution-error", err.Error())
		}
		return false
	}

	// 6. 处理结果
	if cb.EmitEvent != nil {
		cb.EmitEvent("solution", response.Content)
	}

	if req.Config.GetKeepContext() {
		// 保持上下文模式：保存完整的用户消息和助手回复到历史
		s.chatHistory = append(s.chatHistory, currentUserMsg)
		s.chatHistory = append(s.chatHistory, llm.NewAssistantMessage(response.Content))
	} else {
		// 不保持上下文模式：清空历史
		s.chatHistory = []llm.Message{}
	}

	return true
}

// ensureSystemPrompt 确保 chatHistory 的第一条是正确的 System Prompt
func (s *Solver) ensureSystemPrompt(prompt string) {
	if len(s.chatHistory) == 0 {
		s.chatHistory = append(s.chatHistory, llm.NewSystemMessage(prompt))
		logger.Println("插入 SystemPrompt")
		return
	}

	// 检查第一条是否为系统消息
	if s.chatHistory[0].Role == llm.RoleSystem {
		if s.chatHistory[0].Content != prompt {
			s.chatHistory[0] = llm.NewSystemMessage(prompt)
			logger.Println("替换 SystemPrompt")
		}
	} else {
		// 第一条不是系统消息，插入到头部
		s.chatHistory = append([]llm.Message{llm.NewSystemMessage(prompt)}, s.chatHistory...)
		logger.Println("插入 SystemPrompt 到消息历史头部")
	}
}
