package llm

import "context"

// Provider LLM 提供商接口
type Provider interface {
	// GenerateContentStream 流式生成内容
	// messages: 统一格式的消息列表
	// onChunk: 流式输出回调，每个 chunk 包含类型（thinking/content）和内容
	// 返回完整的助手消息（包含 thinking 和 content）
	GenerateContentStream(ctx context.Context, messages []Message, onChunk StreamCallback) (Message, error)

	// GenerateContent 非流式生成内容（可指定模型）
	// model: 模型名称，为空则使用配置中的默认模型
	// messages: 统一格式的消息列表
	// 返回完整的助手消息
	GenerateContent(ctx context.Context, model string, messages []Message) (Message, error)

	// GetModels 获取可用模型列表
	GetModels(ctx context.Context) ([]string, error)

	// TestChat 测试连通性
	TestChat(ctx context.Context) error
}
