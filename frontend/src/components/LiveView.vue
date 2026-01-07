<template>
  <div class="live-view">
    <!-- 顶部状态栏 -->
    <div class="live-header">
      <div class="header-title">
        <span class="live-dot" :class="statusClass"></span>
        <span class="title-text">实时对话</span>
      </div>
      <div class="header-status">{{ statusText }}</div>
    </div>

    <!-- 聊天区域 -->
    <div class="chat-area" ref="chatContainer">
      <!-- 空状态 -->
      <div v-if="messages.length === 0" class="empty-state">
        <div class="empty-visual">
          <div class="pulse-ring"></div>
          <div class="pulse-ring delay"></div>
          <span class="mic-icon">🎤</span>
        </div>
        <div class="empty-title">准备就绪</div>
        <div class="empty-desc">正在监听面试对话...</div>
      </div>

      <!-- 消息列表 -->
      <template v-else>
        <div v-for="msg in messages" :key="msg.id" class="msg-wrapper" :class="msg.type">
          <div class="msg-card" :class="{ interrupted: msg.interrupted }">
            <div class="msg-header">
              <span class="msg-sender">{{ msg.type === 'interviewer' ? '面试官' : 'AI 建议' }}</span>
              <span v-if="msg.interrupted" class="interrupted-tag">已打断</span>
              <span class="msg-time">{{ formatTime(msg.timestamp) }}</span>
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
  padding: 14px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.header-title {
  display: flex;
  align-items: center;
  gap: 10px;
}

.live-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  transition: all 0.3s;
}

.live-dot.disconnected {
  background: #6b7280;
}

.live-dot.connecting {
  background: #f59e0b;
  animation: blink 1s infinite;
}

.live-dot.connected {
  background: #10b981;
  box-shadow: 0 0 8px rgba(16, 185, 129, 0.6);
}

.live-dot.error {
  background: #ef4444;
}

@keyframes blink {

  0%,
  100% {
    opacity: 1;
  }

  50% {
    opacity: 0.3;
  }
}

.title-text {
  font-size: 14px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
  letter-spacing: 0.3px;
}

.header-status {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.45);
}

/* ===== 聊天区域 ===== */
.chat-area {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 16px 20px 60px 20px;
  /* 底部增加更多 padding */
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-height: 0;
  pointer-events: auto;
}

.chat-area::-webkit-scrollbar {
  width: 4px;
}

.chat-area::-webkit-scrollbar-track {
  background: transparent;
}

.chat-area::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.08);
  border-radius: 2px;
}

.scroll-spacer {
  height: 40px;
  /* 增大底部空白 */
  flex-shrink: 0;
}

/* ===== 空状态 ===== */
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
}

.empty-visual {
  position: relative;
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pulse-ring {
  position: absolute;
  width: 100%;
  height: 100%;
  border: 2px solid rgba(16, 185, 129, 0.3);
  border-radius: 50%;
  animation: pulse-out 2s ease-out infinite;
}

.pulse-ring.delay {
  animation-delay: 1s;
}

@keyframes pulse-out {
  0% {
    transform: scale(0.5);
    opacity: 1;
  }

  100% {
    transform: scale(1.5);
    opacity: 0;
  }
}

.mic-icon {
  font-size: 28px;
  z-index: 1;
}

.empty-title {
  font-size: 16px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.8);
}

.empty-desc {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.4);
}

/* ===== 消息卡片 ===== */
.msg-wrapper {
  display: flex;
  animation: fadeSlide 0.25s ease-out;
}

.msg-wrapper.interviewer {
  justify-content: flex-start;
}

.msg-wrapper.ai {
  justify-content: flex-end;
}

@keyframes fadeSlide {
  from {
    opacity: 0;
    transform: translateY(6px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.msg-card {
  max-width: 80%;
  padding: 12px 16px;
  border-radius: 12px;
}

/* 面试官消息 - 优雅的蓝灰渐变 */
.interviewer .msg-card {
  background: linear-gradient(135deg, rgba(55, 65, 81, 0.95) 0%, rgba(45, 55, 72, 0.9) 100%);
  border: 1px solid rgba(148, 163, 184, 0.15);
  border-left: 3px solid rgba(148, 163, 184, 0.4);
  border-radius: 2px 12px 12px 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

/* AI 消息 - 柔和的品牌色 */
.ai .msg-card {
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.18) 0%, rgba(6, 182, 212, 0.12) 100%);
  border: 1px solid rgba(16, 185, 129, 0.25);
  border-right: 3px solid rgba(16, 185, 129, 0.5);
  border-radius: 12px 2px 12px 12px;
  box-shadow: 0 2px 8px rgba(16, 185, 129, 0.1);
}

.msg-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
}

.msg-sender {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.interviewer .msg-sender {
  color: rgba(148, 163, 184, 0.9);
}

.ai .msg-sender {
  color: rgba(16, 185, 129, 0.95);
}

.msg-time {
  font-size: 10px;
  color: rgba(255, 255, 255, 0.3);
  margin-left: auto;
}

/* 打断标签 */
.interrupted-tag {
  font-size: 9px;
  font-weight: 600;
  color: #f59e0b;
  background: rgba(245, 158, 11, 0.15);
  padding: 2px 6px;
  border-radius: 4px;
  margin-left: 6px;
}

/* 被打断的消息卡片 */
.msg-card.interrupted {
  opacity: 0.85;
  border-style: dashed;
}

.msg-body {
  font-size: 14px;
  line-height: 1.55;
  color: rgba(255, 255, 255, 0.92);
}

.msg-body :deep(p) {
  margin: 0 0 6px 0;
}

.msg-body :deep(p:last-child) {
  margin-bottom: 0;
}

.msg-body :deep(code) {
  background: rgba(0, 0, 0, 0.2);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 13px;
  font-family: 'SF Mono', Consolas, monospace;
}

.msg-body :deep(pre) {
  background: rgba(0, 0, 0, 0.25);
  padding: 10px;
  border-radius: 6px;
  margin: 6px 0;
  overflow-x: auto;
}

.msg-body :deep(pre code) {
  background: none;
  padding: 0;
}

/* 打字指示器 */
.typing-dots {
  display: flex;
  gap: 4px;
  margin-top: 8px;
}

.typing-dots span {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.4);
  animation: dotPulse 1.2s infinite ease-in-out;
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

@keyframes dotPulse {

  0%,
  60%,
  100% {
    transform: scale(1);
    opacity: 0.4;
  }

  30% {
    transform: scale(1.3);
    opacity: 1;
  }
}
</style>
