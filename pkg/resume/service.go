package resume

import (
	"Q-Solver/pkg/config"
	"Q-Solver/pkg/llm"
	"Q-Solver/pkg/logger"
	"Q-Solver/pkg/prompts"
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Service struct {
	config       config.Config // 存储配置副本
	resumeBase64 string        // 缓存的简历 Base64
}

func NewService(cfg config.Config, cm *config.ConfigManager) *Service {
	s := &Service{
		config: cfg,
	}
	// 订阅配置变更，同步配置
	cm.Subscribe(func(NewConfig config.Config, oldConfig config.Config) {
		s.config = NewConfig
		// 如果简历路径变了，清空缓存
		if NewConfig.ResumePath != oldConfig.ResumePath {
			s.resumeBase64 = ""
		}
	})
	return s
}

// SelectResume 打开文件对话框选择简历，并返回路径
func (s *Service) SelectResume(ctx context.Context) string {
	selection, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
		Title: "选择简历 (PDF)",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "PDF Files",
				Pattern:     "*.pdf",
			},
		},
	})

	if err != nil {
		logger.Printf("选择文件失败: %v\n", err)
		return ""
	}

	if selection == "" {
		return "" // 用户取消
	}

	// 不再直接修改配置，由调用者处理
	return selection
}

// ClearResume 清除简历缓存
func (s *Service) ClearResume() {
	s.resumeBase64 = ""
	logger.Println("简历缓存已清除")
}

// GetResumeBase64 读取简历并转换为 Base64
func (s *Service) GetResumeBase64() (string, error) {
	// 读一下缓存
	if len(s.resumeBase64) > 0 {
		logger.Println("使用缓存的简历 Base64")
		return s.resumeBase64, nil
	}
	if s.config.ResumePath == "" {
		return "", nil
	}

	// 检查文件大小
	fileInfo, err := os.Stat(s.config.ResumePath)
	if err != nil {
		return "", err
	}

	// 限制 5MB
	if fileInfo.Size() > 5*1024*1024 {
		return "", fmt.Errorf("简历文件大小超过 5MB 限制")
	}

	// 读取文件内容
	content, err := os.ReadFile(s.config.ResumePath)
	if err != nil {
		return "", err
	}

	// 转换为 Base64 并缓存
	encoded := base64.StdEncoding.EncodeToString(content)
	s.resumeBase64 = encoded
	return encoded, nil
}

// ParseResume 解析简历为 Markdown
func (s *Service) ParseResume(ctx context.Context, provider llm.Provider) (string, error) {
	// 1. Read Resume
	resumeBase64, err := s.GetResumeBase64()
	if err != nil {
		return "", fmt.Errorf("读取简历失败: %v", err)
	}
	if resumeBase64 == "" {
		return "", fmt.Errorf("未选择简历文件")
	}
	logger.Println("开始解析简历为 Markdown...")

	// 2. 构建解析简历的消息
	messages := []llm.Message{
		llm.NewUserMessage(prompts.ResumeParsePrompt),
		llm.NewMultiPartMessage(llm.RoleUser, []llm.ContentPart{
			llm.PDFPart(resumeBase64),
		}),
	}

	// 3. 调用 LLM Provider
	result, err := provider.GenerateContentStream(ctx, messages, nil)
	if err != nil {
		return "", err
	}
	return result.Content, nil
}
