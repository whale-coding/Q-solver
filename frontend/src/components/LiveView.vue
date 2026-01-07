<template>
  <div class="live-view">
    <!-- 顶部状态栏 -->
    <div class="live-header">
      <div class="header-title">
        <div class="audio-bars" :class="{ active: status === 'connected' }">
          <span></span><span></span><span></span>
        </div>
        <span class="title-text">实时助手</span>
      </div>
      <div class="header-status" :class="statusClass">
        <span class="status-dot"></span>
        <span>{{ statusText }}</span>
      </div>
    </div>

    <!-- 错误提示 -->
    <div v-if="status === 'error' && errorMsg" class="error-banner">
      <span class="error-icon">⚠️</span>
      <div class="error-content">
        <div class="error-title">连接失败</div>
        <div class="error-message">{{ errorMsg }}</div>
      </div>
      <button class="retry-btn" @click="retryConnection">重试</button>
    </div>

    <!-- 聊天区域 -->
    <div class="chat-area" ref="chatContainer">
      <!-- 空状态 -->
      <div v-if="messages.length === 0" class="empty-state">
        <div class="empty-visual">
          <div class="wave-container">
            <div class="wave"></div>
            <div class="wave"></div>
            <div class="wave"></div>
          </div>
          <span class="mic-icon">🎙️</span>
        </div>
        <div class="empty-title">准备就绪</div>
        <div class="empty-desc">开始说话，AI 将实时响应</div>
      </div>

      <!-- 消息列表 -->
      <template v-else>
        <div v-for="msg in messages" :key="msg.id" class="msg-wrapper" :class="msg.type">
          <div class="msg-card" :class="{ interrupted: msg.interrupted }">
            <div class="msg-header">
              <div class="sender-info">
                <span class="avatar" :class="msg.type">{{ msg.type === 'interviewer' ? 'Q' : 'A' }}</span>
                <span class="msg-sender">{{ msg.type === 'interviewer' ? '语音' : 'AI' }}</span>
              </div>
              <div class="header-right">
                <span v-if="msg.interrupted" class="interrupted-tag">已中断</span>
                <span class="msg-time">{{ formatTime(msg.timestamp) }}</span>
              </div>
            </div>
            <div class="msg-body" v-html="msg.type === 'ai' ? renderMarkdown(msg.content) : escapeHtml(msg.content)">
            </div>
            <div v-if="!msg.isComplete" class="typing-dots">
              <span></span><span></span><span></span>
            </div>
          </div>
        </div>
      </template>

      <div class="scroll-spacer"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { marked } from 'marked'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'
import { StartLiveSession, StopLiveSession } from '../../wailsjs/go/main/App'

const status = ref('disconnected')
const errorMsg = ref('')
const chatContainer = ref(null)
const messages = ref([])
const currentInterviewerMsg = ref(null)
const currentAiMsg = ref(null)

function generateId() {
  return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
}

function createMessage(type) {
  return {
    id: generateId(),
    type,
    content: '',
    timestamp: Date.now(),
    isComplete: false
  }
}

function formatTime(timestamp) {
  return new Date(timestamp).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

function escapeHtml(text) {
  if (!text) return ''
  return text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/\n/g, '<br>')
}

function renderMarkdown(text) {
  if (!text) return ''
  return marked.parse(text.replace(/\n+$/, ''))
}

const statusClass = computed(() => status.value)
const statusText = computed(() => {
  const map = { disconnected: '等待连接', connecting: '连接中...', connected: '已连接', error: '连接失败' }
  return map[status.value] || '未知'
})

function scrollToBottom() {
  setTimeout(() => {
    if (chatContainer.value) {
      chatContainer.value.scrollTop = chatContainer.value.scrollHeight
    }
  }, 20)
}

watch(messages, scrollToBottom, { deep: true })

function onLiveStatus(s) { status.value = s }
function onLiveTranscript(text) {
  if (!currentInterviewerMsg.value) {
    currentInterviewerMsg.value = createMessage('interviewer')
    messages.value.push(currentInterviewerMsg.value)
  }
  currentInterviewerMsg.value.content += text
}
function onLiveInterviewerDone() {
  if (currentInterviewerMsg.value) {
    currentInterviewerMsg.value.isComplete = true
    currentInterviewerMsg.value = null
  }
}
function onLiveAiText(text) {
  if (!currentAiMsg.value) {
    currentAiMsg.value = createMessage('ai')
    messages.value.push(currentAiMsg.value)
  }
  currentAiMsg.value.content += text
}
function onLiveError(err) { status.value = 'error'; errorMsg.value = err }
function onLiveDone() {
  if (currentAiMsg.value) {
    currentAiMsg.value.isComplete = true
    currentAiMsg.value = null
  }
}

function retryConnection() {
  errorMsg.value = ''
  status.value = 'connecting'
  StopLiveSession()
  StartLiveSession()
}
function onLiveInterrupted(text) {
  // AI 回复被打断，标记当前 AI 消息为已打断
  if (currentAiMsg.value) {
    currentAiMsg.value.isComplete = true
    currentAiMsg.value.interrupted = true
    currentAiMsg.value = null
  }
}

onMounted(() => {
  EventsOn('live:status', onLiveStatus)
  EventsOn('live:transcript', onLiveTranscript)
  EventsOn('live:interviewer-done', onLiveInterviewerDone)
  EventsOn('live:ai-text', onLiveAiText)
  EventsOn('live:error', onLiveError)
  EventsOn('live:done', onLiveDone)
  EventsOn('live:Interrupted', onLiveInterrupted)
  StartLiveSession()
})

onUnmounted(() => {
  StopLiveSession()
  EventsOff('live:status')
  EventsOff('live:transcript')
  EventsOff('live:interviewer-done')
  EventsOff('live:ai-text')
  EventsOff('live:error')
  EventsOff('live:done')
  EventsOff('live:Interrupted')
})
</script>

<style scoped>
/* ===== 基础布局 ===== */
.live-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  pointer-events: auto;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

/* ===== 顶部栏 ===== */
.live-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  position: relative;
}

/* 底部渐变分隔线 - 优雅的中间渐变效果 */
.live-header::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 24px;
  right: 24px;
  height: 1px;
  background: linear-gradient(90deg,
      transparent 0%,
      rgba(16, 185, 129, 0.35) 30%,
      rgba(16, 185, 129, 0.35) 70%,
      transparent 100%);
}

.header-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* 音频条动画 */
.audio-bars {
  display: flex;
  align-items: flex-end;
  gap: 3px;
  height: 16px;
}

.audio-bars span {
  width: 3px;
  height: 4px;
  background: rgba(255, 255, 255, 0.3);
  border-radius: 2px;
  transition: all 0.3s ease;
}

.audio-bars.active span {
  background: linear-gradient(180deg, #10b981 0%, #34d399 100%);
  animation: audioWave 0.8s ease-in-out infinite;
}

.audio-bars.active span:nth-child(1) {
  animation-delay: 0s;
}

.audio-bars.active span:nth-child(2) {
  animation-delay: 0.2s;
}

.audio-bars.active span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes audioWave {

  0%,
  100% {
    height: 4px;
  }

  50% {
    height: 14px;
  }
}

.title-text {
  font-size: 15px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
  letter-spacing: 0.5px;
}

.header-status {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.45);
  transition: color 0.3s ease;
}

.status-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.25);
  transition: all 0.4s ease;
}

.header-status.connecting .status-dot {
  background: #fbbf24;
  box-shadow: 0 0 6px rgba(251, 191, 36, 0.6);
  animation: pulse 1.2s infinite;
}

.header-status.connected {
  color: rgba(16, 185, 129, 0.75);
}

.header-status.connected .status-dot {
  background: #10b981;
  box-shadow: 0 0 6px rgba(16, 185, 129, 0.6);
}

.header-status.error .status-dot {
  background: #ef4444;
  box-shadow: 0 0 8px rgba(239, 68, 68, 0.5);
}

/* ===== 错误提示 ===== */
.error-banner {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 0 16px 12px 16px;
  padding: 14px 16px;
  background: linear-gradient(135deg, rgba(239, 68, 68, 0.15), rgba(185, 28, 28, 0.1));
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 12px;
  backdrop-filter: blur(8px);
  animation: slideDown 0.3s ease-out;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.error-icon {
  font-size: 20px;
  flex-shrink: 0;
}

.error-content {
  flex: 1;
  min-width: 0;
}

.error-title {
  font-size: 13px;
  font-weight: 600;
  color: #fca5a5;
  margin-bottom: 4px;
}

.error-message {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.6);
  word-break: break-word;
  line-height: 1.4;
}

.retry-btn {
  flex-shrink: 0;
  padding: 6px 14px;
  background: linear-gradient(135deg, rgba(239, 68, 68, 0.3), rgba(239, 68, 68, 0.2));
  border: 1px solid rgba(239, 68, 68, 0.4);
  border-radius: 8px;
  color: #fca5a5;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.retry-btn:hover {
  background: linear-gradient(135deg, rgba(239, 68, 68, 0.4), rgba(239, 68, 68, 0.3));
  border-color: rgba(239, 68, 68, 0.6);
  transform: translateY(-1px);
}

@keyframes pulse {

  0%,
  100% {
    opacity: 1;
    transform: scale(1);
  }

  50% {
    opacity: 0.6;
    transform: scale(1.2);
  }
}

/* ===== 聊天区域 ===== */
.chat-area {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 20px 24px 80px 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-height: 0;
  pointer-events: auto;
}

.chat-area::-webkit-scrollbar {
  width: 5px;
}

.chat-area::-webkit-scrollbar-track {
  background: transparent;
}

.chat-area::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 3px;
  transition: background 0.3s;
}

.chat-area:hover::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.15);
}

.scroll-spacer {
  height: 40px;
  flex-shrink: 0;
}

/* ===== 空状态 ===== */
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 20px;
}

.empty-visual {
  position: relative;
  width: 100px;
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.wave-container {
  position: absolute;
  width: 100%;
  height: 100%;
}

.wave {
  position: absolute;
  width: 100%;
  height: 100%;
  border: 2px solid rgba(16, 185, 129, 0.2);
  border-radius: 50%;
  animation: waveExpand 2.5s ease-out infinite;
}

.wave:nth-child(2) {
  animation-delay: 0.8s;
}

.wave:nth-child(3) {
  animation-delay: 1.6s;
}

@keyframes waveExpand {
  0% {
    transform: scale(0.5);
    opacity: 0.8;
  }

  100% {
    transform: scale(1.4);
    opacity: 0;
  }
}

.mic-icon {
  font-size: 32px;
  z-index: 1;
  filter: drop-shadow(0 4px 12px rgba(16, 185, 129, 0.3));
}

.empty-title {
  font-size: 18px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.85);
  letter-spacing: 0.3px;
}

.empty-desc {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.45);
}

/* ===== 消息卡片 ===== */
.msg-wrapper {
  display: flex;
  animation: slideIn 0.3s ease-out;
}

.msg-wrapper.interviewer {
  justify-content: flex-start;
}

.msg-wrapper.ai {
  justify-content: flex-end;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.msg-card {
  max-width: 85%;
  padding: 14px 18px;
  border-radius: 18px;
  backdrop-filter: blur(12px);
  transition: all 0.3s ease;
}

/* 语音消息 - Q */
.interviewer .msg-card {
  background: linear-gradient(135deg, rgba(70, 80, 100, 0.7) 0%, rgba(55, 65, 85, 0.65) 100%);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-left: 3px solid rgba(139, 92, 246, 0.6);
  border-radius: 4px 18px 18px 18px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

/* AI 消息 - A */
.ai .msg-card {
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.2) 0%, rgba(6, 182, 212, 0.15) 100%);
  border: 1px solid rgba(16, 185, 129, 0.2);
  border-right: 3px solid rgba(16, 185, 129, 0.6);
  border-radius: 18px 4px 18px 18px;
  box-shadow: 0 4px 16px rgba(16, 185, 129, 0.1);
}

.msg-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.sender-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

/* 头像样式 */
.avatar {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 0;
}

.avatar.interviewer {
  background: linear-gradient(135deg, #8b5cf6 0%, #a78bfa 100%);
  color: white;
  box-shadow: 0 2px 8px rgba(139, 92, 246, 0.35);
}

.avatar.ai {
  background: linear-gradient(135deg, #10b981 0%, #34d399 100%);
  color: white;
  box-shadow: 0 2px 8px rgba(16, 185, 129, 0.35);
}

.msg-sender {
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.8px;
}

.interviewer .msg-sender {
  color: rgba(167, 139, 250, 0.95);
}

.ai .msg-sender {
  color: rgba(52, 211, 153, 0.95);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.msg-time {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.35);
}

/* 打断标签 */
.interrupted-tag {
  font-size: 10px;
  font-weight: 600;
  color: #fbbf24;
  background: rgba(251, 191, 36, 0.15);
  padding: 3px 8px;
  border-radius: 6px;
  border: 1px solid rgba(251, 191, 36, 0.25);
}

/* 被打断的消息卡片 */
.msg-card.interrupted {
  opacity: 0.75;
  border-style: dashed;
}

.msg-body {
  font-size: 14px;
  line-height: 1.6;
  color: rgba(255, 255, 255, 0.9);
}

.msg-body :deep(p) {
  margin: 0 0 8px 0;
}

.msg-body :deep(p:last-child) {
  margin-bottom: 0;
}

.msg-body :deep(code) {
  background: rgba(0, 0, 0, 0.25);
  padding: 2px 7px;
  border-radius: 5px;
  font-size: 13px;
  font-family: 'SF Mono', Consolas, monospace;
  color: #a5f3fc;
}

.msg-body :deep(pre) {
  background: rgba(0, 0, 0, 0.3);
  padding: 12px;
  border-radius: 10px;
  margin: 8px 0;
  overflow-x: auto;
  border: 1px solid rgba(255, 255, 255, 0.06);
}

.msg-body :deep(pre code) {
  background: none;
  padding: 0;
  color: #e2e8f0;
}

/* 打字指示器 */
.typing-dots {
  display: flex;
  gap: 5px;
  margin-top: 10px;
  padding-top: 8px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.typing-dots span {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.4);
  animation: bounce 1.4s infinite ease-in-out;
}

.typing-dots span:nth-child(1) {
  animation-delay: 0s;
}

.typing-dots span:nth-child(2) {
  animation-delay: 0.15s;
}

.typing-dots span:nth-child(3) {
  animation-delay: 0.3s;
}

@keyframes bounce {

  0%,
  60%,
  100% {
    transform: translateY(0);
    opacity: 0.4;
  }

  30% {
    transform: translateY(-6px);
    opacity: 1;
  }
}
</style>
