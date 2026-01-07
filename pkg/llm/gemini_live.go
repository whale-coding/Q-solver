package llm

import (
	"Q-Solver/pkg/config"
	"Q-Solver/pkg/logger"
	"context"

	"google.golang.org/genai"
)

// GeminiLiveSession 封装 Gemini SDK 的 Live 会话
type GeminiLiveSession struct {
	session             *genai.Session
	hasInterviewerInput bool         // 是否有面试官输入待处理
	pendingAIMessage    *LiveMessage // 缓存的 AI 消息（等待 InterviewerDone 发送后再返回）
}

// ConnectLive 实现 LiveProvider 接口
func (a *GeminiAdapter) ConnectLive(ctx context.Context, cfg *LiveConfig, config *config.Config) (LiveSession, error) {
	// Live API 只支持特定模型，用户配置的模型可能不支持 bidiGenerateContent
	model := cfg.Model
	if model == "" {
		model = a.config.GetModel()
	}
	// 定义截图工具
	screenshotTool := &genai.Tool{
		FunctionDeclarations: []*genai.FunctionDeclaration{
			{
				Name:        "get_screenshot",
				Description: "获取用户当前屏幕截图，用于查看题目或界面内容",
			},
		},
	}

	// 连接配置
	connectCfg := &genai.LiveConnectConfig{
		Tools:                   []*genai.Tool{screenshotTool},
		ResponseModalities:      []genai.Modality{genai.ModalityAudio},
		MaxOutputTokens:         int32(config.GetMaxTokens()),
		Temperature:             toFloat32Ptr(config.GetTemperature()),
		TopP:                    toFloat32Ptr(config.GetTopP()),
		TopK:                    intToFloat32Ptr(config.GetTopK()),
		InputAudioTranscription: &genai.AudioTranscriptionConfig{},
		SpeechConfig: &genai.SpeechConfig{
			LanguageCode: "cmn-CN",
			VoiceConfig: &genai.VoiceConfig{
				PrebuiltVoiceConfig: &genai.PrebuiltVoiceConfig{
					VoiceName: "Aoede",
				},
			},
		},
		OutputAudioTranscription: &genai.AudioTranscriptionConfig{},
	}

	// 定义结构化系统指令
	const InterviewCopilotInstruction = `
# 角色与人设
你是一个名为“Q-Solver”的精英级技术面试助手（Copilot）。你的语气应当专业、简洁且充满支持性。请注意：你是在直接辅助我（候选人），而不是在与面试官对话。

# 主要目标
实时聆听面试语境，并提供即时、高质量的答案或谈话要点，以便我能直接复述给面试官。

# 语言强制要求
- **输出语言**：无论面试官使用的是何种语言（如英文），无论是你的回复（包括答案、建议、解释），还是你对用户语音的转录识别，必须严格使用**简体中文**。

# 对话规则
1. **等待提问（至关重要）**：
   - 严禁在面试官解释背景、陈述上下文或刚开始说话时插嘴。
   - **必须**等待听到明确的问题，或出现表明轮到我说话的明显停顿。

2. **顾问模式**：
   - 不要试图与面试官对话。请直接对我说话。
   - 回答时请以“你可以这样说……”开头，或者直接列出简洁的要点。

3. **视觉感知**：
   - 如果面试官提到了屏幕上的代码、架构图或文本，而你无法通过音频获取信息，请**立即**调用 get_screenshot 工具来分析视觉内容。

4. **回答风格**：
   - **结构化**：采用“第一点、第二点、结论”的结构。
   - **关键词**：使用技术面试中高频的得分关键词（英文术语可保留，但解释用中文）。
   - **极度简洁**：我需要一边听你说的内容一边复述，所以请不要长篇大论，只给核心干货。

# 安全护栏
- 如果你不确定答案，请提供一个我可以反问面试官的“澄清问题”。
- 不要凭空臆造（Hallucinate）需求。如果音频不清晰，直接告诉我“我没听清”。
- **严重警告**：绝对不要打断对话。如果你不确定面试官是否说完，宁愿保持沉默也不要抢话。
`

	instructionText := cfg.SystemInstruction
	if instructionText == "" {
		instructionText = InterviewCopilotInstruction
	} else {
		// 如果用户有自定义指令，将其合并以保留核心规则
		instructionText = InterviewCopilotInstruction + "\n\n# User Preferences\n" + instructionText
	}

	connectCfg.SystemInstruction = &genai.Content{
		Parts: []*genai.Part{{Text: instructionText}},
	}

	session, err := a.client.Live.Connect(ctx, model, connectCfg)
	if err != nil {
		logger.Printf("LiveAPI: 连接到模型 %s 发生错误", err)
		return nil, err
	}
	logger.Printf("LiveAPI: 连接到模型 %s", model)
	return &GeminiLiveSession{session: session}, nil
}

// SendAudio 发送音频数据 (16kHz, 16-bit, mono PCM)
func (s *GeminiLiveSession) SendAudio(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	return s.session.SendRealtimeInput(genai.LiveRealtimeInput{
		Media: &genai.Blob{
			MIMEType: "audio/pcm;rate=16000",
			Data:     data,
		},
	})
}

// Receive 接收消息 (阻塞)
func (s *GeminiLiveSession) Receive() (*LiveMessage, error) {
	// 如果有缓存的 AI 消息，先返回它
	if s.pendingAIMessage != nil {
		msg := s.pendingAIMessage
		s.pendingAIMessage = nil
		return msg, nil
	}

	msg, err := s.session.Receive()
	if err != nil {
		return &LiveMessage{Type: LiveMsgError, Text: err.Error()}, err
	}
	return s.convertMessage(msg), nil
}

// convertMessage 转换 SDK 消息为统一格式
func (s *GeminiLiveSession) convertMessage(msg *genai.LiveServerMessage) *LiveMessage {
	if msg == nil {
		return nil
	}
	if msg.ServerContent != nil && msg.ServerContent.Interrupted{
		return  &LiveMessage{Type: LiveInterrupted, Text: "检测到面试官说话(已打断当前回复)"}
	}
	// 输入音频转录 (面试官说话的文字)
	if msg.ServerContent != nil && msg.ServerContent.InputTranscription != nil {
		text := msg.ServerContent.InputTranscription.Text
		if text != "" {
			s.hasInterviewerInput = true // 标记有面试官输入
			return &LiveMessage{Type: LiveMsgTranscript, Text: text}
		}
	}

	if msg.ServerContent != nil && msg.ServerContent.OutputTranscription != nil {
		text := msg.ServerContent.OutputTranscription.Text
		if text != "" {
			aiMsg := &LiveMessage{Type: LiveMsgAIText, Text: text}
			// 如果之前有面试官输入，先发送面试官结束信号，缓存 AI 消息
			if s.hasInterviewerInput {
				s.hasInterviewerInput = false
				s.pendingAIMessage = aiMsg // 缓存 AI 消息，下次 Receive 时返回
				return &LiveMessage{Type: LiveMsgInterviewerDone}
			}
			return aiMsg
		}
	}

	// 工具调用
	if msg.ToolCall != nil && len(msg.ToolCall.FunctionCalls) > 0 {
		fc := msg.ToolCall.FunctionCalls[0]
		return &LiveMessage{
			Type:     LiveMsgToolCall,
			ToolName: fc.Name,
			ToolID:   fc.ID,
		}
	}

	// 服务端消息
	if msg.ServerContent != nil {
		// 是否完成
		if msg.ServerContent.TurnComplete {
			return &LiveMessage{Type: LiveMsgDone}
		}

		// 检查 ModelTurn 中的文本 (当 ResponseModalities 为 Text 时)
		if msg.ServerContent.ModelTurn != nil {
			for _, part := range msg.ServerContent.ModelTurn.Parts {
				if part != nil && part.Text != "" {
					return &LiveMessage{Type: LiveMsgAIText, Text: part.Text}
				}
			}
		}
	}

	return nil
}

// SendToolResponse 发送工具调用结果 (文本)
func (s *GeminiLiveSession) SendToolResponse(toolID string, result string) error {
	return s.session.SendToolResponse(genai.LiveToolResponseInput{
		FunctionResponses: []*genai.FunctionResponse{
			{
				ID:       toolID,
				Response: map[string]any{"content": result},
			},
		},
	})
}

// SendToolResponseWithImage 发送图片作为工具调用结果
func (s *GeminiLiveSession) SendToolResponseWithImage(toolID string, imageData []byte, mimeType string) error {
	logger.Printf("LiveAPI: 发送图片工具响应 ID=%s, size=%d, mime=%s", toolID, len(imageData), mimeType)
	return s.session.SendToolResponse(genai.LiveToolResponseInput{
		FunctionResponses: []*genai.FunctionResponse{
			{
				ID: toolID,
				Response: map[string]any{
					"image": map[string]any{
						"mimeType": mimeType,
						"data":     imageData,
					},
				},
			},
		},
	})
}

// Close 关闭会话
func (s *GeminiLiveSession) Close() error {
	return s.session.Close()
}

// SupportsLive 检查是否支持 Live API
func SupportsLive(p Provider) bool {
	_, ok := p.(LiveProvider)
	return ok
}

// GetLiveConfig 从配置创建 LiveConfig
func GetLiveConfig(cfg *config.Config) *LiveConfig {
	return &LiveConfig{
		Model:             cfg.GetModel(),
		SystemInstruction: cfg.GetPrompt(),
	}
}

func toFloat32Ptr(v float64) *float32 {
	f := float32(v)
	return &f
}

func intToFloat32Ptr(v int) *float32 {
	f := float32(v)
	return &f
}
