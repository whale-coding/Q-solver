package live

import (
	"Q-Solver/pkg/audio"
	"Q-Solver/pkg/config"
	"Q-Solver/pkg/llm"
	"Q-Solver/pkg/logger"
	"Q-Solver/pkg/screen"
	"context"
	"encoding/base64"
	"strings"
	"sync"
)

// LiveSessionManager 管理 Live API 会话的完整生命周期
type LiveSessionManager struct {
	// 依赖注入
	ctx           context.Context
	llmService    *llm.Service
	configManager *config.ConfigManager
	screenService *screen.Service
	emitEvent     func(string, ...any)

	// Live Session 状态
	session      llm.LiveSession
	audioCapture *audio.LoopbackCapture
	mu           sync.Mutex

	// 协程管理
	stopChan  chan struct{}
	errorChan chan error // 错误通道，用于协程报告异常
	wg        sync.WaitGroup
}

// NewLiveSessionManager 创建 Live Session 管理器
func NewLiveSessionManager(
	ctx context.Context,
	llmService *llm.Service,
	configManager *config.ConfigManager,
	screenService *screen.Service,
	emitEvent func(string, ...any),
) *LiveSessionManager {
	return &LiveSessionManager{
		ctx:           ctx,
		llmService:    llmService,
		configManager: configManager,
		screenService: screenService,
		emitEvent:     emitEvent,
	}
}

// Start 启动 Live API 会话
func (m *LiveSessionManager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	cfg := m.configManager.Get()

	// 检查 Provider 是否支持 Live
	provider := m.llmService.GetProvider()
	liveProvider, ok := provider.(llm.LiveProvider)
	if !ok {
		m.emitEvent("live:error", "当前模型不支持 Live API")
		return nil
	}

	m.emitEvent("live:status", "connecting")

	// 连接 Live Session
	liveCfg := llm.GetLiveConfig(cfg)
	session, err := liveProvider.ConnectLive(m.ctx, liveCfg)
	if err != nil {
		logger.Println("liveApi连接服务器失败", err)
		m.emitEvent("live:status", "error")
		m.emitEvent("live:error", err.Error())
		return err
	}

	m.emitEvent("live:status", "connected")

	// 保存 session
	m.session = session

	// 初始化音频采集
	m.audioCapture, err = audio.NewLoopbackCapture(nil)
	if err != nil {
		logger.Printf("音频采集初始化失败: %v", err)
		m.emitEvent("live:error", "音频采集初始化失败: "+err.Error())
		// 音频采集失败，关闭会话
		session.Close()
		m.session = nil
		m.emitEvent("live:status", "error")
		return err
	}

	if err := m.audioCapture.Start(); err != nil {
		logger.Printf("音频采集启动失败: %v", err)
		m.emitEvent("live:error", "音频采集启动失败: "+err.Error())
		// 音频采集失败，关闭会话和音频设备
		m.audioCapture.Close()
		m.audioCapture = nil
		session.Close()
		m.session = nil
		m.emitEvent("live:status", "error")
		return err
	}

	// 初始化通道
	m.stopChan = make(chan struct{})
	m.errorChan = make(chan error, 2) // 缓冲区为2，防止阻塞

	// 启动错误监听协程
	go m.errorWatcher()

	// 启动音频发送协程
	m.wg.Add(1)
	go m.audioSender(session, m.audioCapture.GetAudioChannel())

	// 启动接收协程
	m.wg.Add(1)
	go m.receiveLoop(session)

	return nil
}

// Stop 停止 Live API 会话（外部调用）
func (m *LiveSessionManager) Stop() {
	// 发送停止信号，通知 errorWatcher 退出
	m.mu.Lock()
	if m.stopChan != nil {
		select {
		case <-m.stopChan:
			// 已关闭
		default:
			close(m.stopChan)
		}
	}
	m.mu.Unlock()

	// 执行清理
	m.cleanup()

	// 等待协程结束
	m.wg.Wait()

	m.emitEvent("live:status", "disconnected")
}

// cleanup 内部清理方法（停止音频采集和关闭会话）
func (m *LiveSessionManager) cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	logger.Println("Live: cleanup 调用")

	// 停止音频采集
	if m.audioCapture != nil {
		logger.Println("Live: 停止音频采集")
		m.audioCapture.Close()
		m.audioCapture = nil
	}

	// 关闭会话
	if m.session != nil {
		logger.Println("Live: 关闭会话")
		m.session.Close()
		m.session = nil
	}
}

// errorWatcher 监听错误通道，发生异常时执行清理
func (m *LiveSessionManager) errorWatcher() {
	select {
	case err := <-m.errorChan:
		logger.Printf("Live: errorWatcher 收到错误: %v", err)
		// 执行清理
		m.cleanup()
		// 通知前端
		m.emitEvent("live:status", "error")
		m.emitEvent("live:error", err.Error())
	case <-m.stopChan:
		// 正常停止，不做额外处理
		logger.Println("Live: errorWatcher 正常退出")
	}
}

// IsActive 检查 Live Session 是否活跃
func (m *LiveSessionManager) IsActive() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.session != nil
}

// audioSender 从音频 channel 读取数据并发送给 Live Session
func (m *LiveSessionManager) audioSender(session llm.LiveSession, audioChan <-chan []byte) {
	defer m.wg.Done()

	logger.Println("Live: 音频发送协程已启动")
	for audioData := range audioChan {
		if session != nil {
			if err := session.SendAudio(audioData); err != nil {
				logger.Printf("Live: 发送音频失败: %v", err)
				// 向错误通道报告异常
				select {
				case m.errorChan <- err:
				default:
					// 通道满或已关闭，忽略
				}
				return
			}
		}
	}
	logger.Println("Live: 音频发送协程已结束")
}

// receiveLoop 接收 Live 消息的循环
func (m *LiveSessionManager) receiveLoop(session llm.LiveSession) {
	defer m.wg.Done()

	logger.Println("Live: 接收循环已启动")

	for {
		msg, err := session.Receive()
		if err != nil {
			logger.Printf("Live: 接收错误: %v", err)
			// 向错误通道报告异常
			select {
			case m.errorChan <- err:
			default:
				// 通道满或已关闭，忽略
			}
			return
		}
		if msg == nil {
			continue
		}

		switch msg.Type {
		case llm.LiveInterrupted:
			logger.Println("检测到打断")
			m.emitEvent("live:Interrupted", msg.Text)
		case llm.LiveMsgTranscript:
			m.emitEvent("live:transcript", msg.Text)
		case llm.LiveMsgInterviewerDone:
			logger.Println("Live: 面试官说话结束")
			m.emitEvent("live:interviewer-done")
		case llm.LiveMsgAIText:
			m.emitEvent("live:ai-text", msg.Text)
		case llm.LiveMsgToolCall:
			logger.Printf("Live: 工具调用 %s (ID=%s)", msg.ToolName, msg.ToolID)
			if msg.ToolName == "get_screenshot" {
				m.handleScreenshot(session, msg.ToolID)
			}
		case llm.LiveMsgDone:
			logger.Println("Live: 对话轮完成")
			m.emitEvent("live:done")
		case llm.LiveMsgError:
			logger.Printf("Live: 错误: %s", msg.Text)
			m.emitEvent("live:error", msg.Text)
		}
	}
}

// handleScreenshot 处理 Live API 的截图请求
func (m *LiveSessionManager) handleScreenshot(session llm.LiveSession, toolID string) {
	cfg := m.configManager.Get()

	preview, err := m.screenService.CapturePreview(
		cfg.CompressionQuality,
		cfg.Sharpening,
		cfg.Grayscale,
		cfg.NoCompression,
		cfg.ScreenshotMode,
	)
	if err != nil {
		logger.Printf("Live 截图失败: %v", err)
		_ = session.SendToolResponse(toolID, "截图失败: "+err.Error())
		return
	}

	// 解析 data URL 格式: data:image/jpeg;base64,xxxxx
	base64Str := preview.Base64
	mimeType := "image/jpeg" // 默认

	if strings.HasPrefix(base64Str, "data:") {
		// 提取 MIME 类型
		if idx := strings.Index(base64Str, ";base64,"); idx > 5 {
			mimeType = base64Str[5:idx]   // 提取 "image/jpeg" 或 "image/png"
			base64Str = base64Str[idx+8:] // 去掉前缀，保留纯 base64
		}
	}

	// 解码 Base64 为原始图片数据
	imageData, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		logger.Printf("Live Base64解码失败: %v", err)
		_ = session.SendToolResponse(toolID, "图片解码失败")
		return
	}

	// 发送图片数据给模型
	err = session.SendToolResponseWithImage(toolID, imageData, mimeType)
	if err != nil {
		logger.Printf("Live 发送截图失败: %v", err)
	} else {
		logger.Printf("Live: 已发送屏幕截图给模型 (%d bytes, %s)", len(imageData), mimeType)
	}
}
