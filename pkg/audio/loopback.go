package audio

import (
	"Q-Solver/pkg/common"
	"Q-Solver/pkg/logger"
	"sync"
	"time"

	"github.com/gen2brain/malgo"
)

const (
	PacketDurationMs = 30                                                    // 每个数据包的时长（毫秒）
	SampleRate       = 16000                                                 // 采样率 16kHz
	BytesPerSample   = 2                                                     // S16 格式，每个样本 2 字节
	PacketSize       = SampleRate * BytesPerSample * PacketDurationMs / 1000 // 每包字节数
	ChannelCapacity  = 100                                                   // channel 容量（可存储的数据包数量）
	RingBufferSize   = PacketSize * 200                                      // 环形缓冲区大小（可存储 200 个数据包，约 8 秒）
)

// LoopbackCapture 扬声器音频采集（使用环形缓冲区 + channel）
type LoopbackCapture struct {
	ctx     *malgo.AllocatedContext
	device  *malgo.Device
	mu      sync.Mutex
	running bool

	// 环形缓冲区：音频采集回调直接写入
	ringBuffer *common.RingBuffer

	// channel：消费者从此读取固定大小的音频包
	audioChan chan []byte

	// 停止信号
	stopChan chan struct{}
	wg       sync.WaitGroup
}

// NewLoopbackCapture 创建 Loopback 采集器
func NewLoopbackCapture(onData func([]byte)) (*LoopbackCapture, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, err
	}
	return &LoopbackCapture{
		ctx:        ctx,
		ringBuffer: common.NewRingBuffer(RingBufferSize),
		audioChan:  make(chan []byte, ChannelCapacity),
		stopChan:   make(chan struct{}),
	}, nil
}

// Start 开始采集扬声器输出
func (c *LoopbackCapture) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return nil
	}

	// 配置音频设备
	deviceConfig := malgo.DefaultDeviceConfig(malgo.Loopback)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.SampleRate = SampleRate

	// 数据回调 - 直接写入环形缓冲区
	onRecv := func(_, pInput []byte, frameCount uint32) {
		if len(pInput) > 0 {
			_, _ = c.ringBuffer.Write(pInput)
		}
	}

	callbacks := malgo.DeviceCallbacks{Data: onRecv}

	device, err := malgo.InitDevice(c.ctx.Context, deviceConfig, callbacks)
	if err != nil {
		return err
	}

	if err := device.Start(); err != nil {
		device.Uninit()
		return err
	}

	c.device = device
	c.running = true

	// 启动消费者协程：从环形缓冲区读取固定大小数据包并发送到 channel
	c.wg.Add(1)
	go c.packetizer()

	logger.Println("Loopback 采集已启动（环形缓冲区模式）")
	return nil
}

// packetizer 从环形缓冲区读取固定大小的数据包并发送到 channel
func (c *LoopbackCapture) packetizer() {
	defer c.wg.Done()

	ticker := time.NewTicker(time.Duration(PacketDurationMs) * time.Millisecond)
	defer ticker.Stop()

	packetBuffer := make([]byte, PacketSize)

	for {
		select {
		case <-c.stopChan:
			logger.Println("音频分包器已停止")
			return

		case <-ticker.C:
			// 尝试从环形缓冲区读取一个完整的数据包
			n, err := c.ringBuffer.Read(packetBuffer)
			if err == common.ErrNotEnoughData {
				// 缓冲区数据不足，跳过这次
				continue
			}
			if err != nil {
				logger.Printf("读取环形缓冲区失败: %v", err)
				continue
			}

			if n == PacketSize {
				// 创建副本发送到 channel
				packet := make([]byte, PacketSize)
				copy(packet, packetBuffer)

				select {
				case c.audioChan <- packet:
					// 成功发送
				default:
					// channel 满了，丢弃此包（或选择覆盖最旧的）
					logger.Println("音频 channel 已满，丢弃数据包")
				}
			}
		}
	}
}

// GetAudioChannel 获取音频数据 channel（供 Live API 消费）
func (c *LoopbackCapture) GetAudioChannel() <-chan []byte {
	return c.audioChan
}

// Stop 停止采集
func (c *LoopbackCapture) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running {
		return
	}

	// 停止设备
	if c.device != nil {
		c.device.Stop()
		c.device.Uninit()
		c.device = nil
	}

	// 停止分包器协程
	close(c.stopChan)
	c.wg.Wait()

	// 重置环形缓冲区
	c.ringBuffer.Reset()

	// 清空 channel
	for len(c.audioChan) > 0 {
		<-c.audioChan
	}

	c.running = false
	logger.Println("Loopback 采集已停止")
}

// Close 释放资源
func (c *LoopbackCapture) Close() {
	c.Stop()

	// 关闭 channel
	close(c.audioChan)

	if c.ctx != nil {
		_ = c.ctx.Uninit()
		c.ctx.Free()
		c.ctx = nil
	}
}

// IsRunning 是否正在采集
func (c *LoopbackCapture) IsRunning() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.running
}

// GetBufferStatus 获取缓冲区状态（用于监控）
func (c *LoopbackCapture) GetBufferStatus() (bufferSize int, channelSize int) {
	return c.ringBuffer.Len(), len(c.audioChan)
}
