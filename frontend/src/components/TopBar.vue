<template>
  <div class="top-bar-wrapper" style="--wails-draggable:drag">
    <div class="top-bar">
      <div class="control-group" :class="{ active: activeButtons.toggle }" style="--wails-draggable:no-drag">
        <span class="key-hint">{{ shortcuts.toggle?.keyName || (isMacOS ? '⌘2' : 'F9') }}</span>
        <span class="label">隐藏/展示</span>
      </div>
      <div class="control-group" :class="{ active: activeButtons.solve }" style="--wails-draggable:no-drag">
        <span class="key-hint">{{ shortcuts.solve?.keyName || (isMacOS ? '⌘1' : 'F8') }}</span>
        <span class="label">一键解题</span>
      </div>
      <div class="control-group" :class="{ active: activeButtons.clickthrough || isClickThrough }" style="--wails-draggable:no-drag">
        <span class="key-hint">{{ shortcuts.clickthrough?.keyName || (isMacOS ? '⌘3' : 'F10') }}</span>
        <span class="label">鼠标穿透</span>
      </div>
      <div class="control-group" style="cursor: default;">
        <span class="key-hint">{{ isMacOS ? '⌘⌥+Move' : 'Alt+Move' }}</span>
        <span class="label">移动/滚动</span>
      </div>
      <div class="divider"></div>
      <div class="control-group" @click="$emit('openSettings')" style="cursor: pointer; --wails-draggable:no-drag"
        @mouseenter="showSettingsTooltip" @mouseleave="hideSettingsTooltip" ref="settingsBtnRef">
        <span class="label">⚙️ 设置</span>
      </div>
      <div class="divider"></div>
      <div class="status-group" ref="statusGroupRef" @mouseenter="showTooltip" @mouseleave="hideTooltip" style="--wails-draggable:no-drag">
        <div class="status-indicator" :class="statusClass">
          <svg class="status-svg" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <!-- 已连接/就绪/解题完成: 信号满格 -->
            <template v-if="isConnectedStatus">
              <rect x="2" y="16" width="4" height="6" rx="1" fill="currentColor"/>
              <rect x="8" y="11" width="4" height="11" rx="1" fill="currentColor"/>
              <rect x="14" y="6" width="4" height="16" rx="1" fill="currentColor"/>
              <rect x="20" y="2" width="2" height="20" rx="1" fill="currentColor"/>
            </template>
            <!-- 未配置: 齿轮 -->
            <template v-else-if="isUnconfigured">
              <circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2"/>
              <path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            </template>
            <!-- Key无效: 警告 -->
            <template v-else-if="isInvalidKey">
              <path d="M12 9v4M12 17h.01" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
              <path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
            </template>
            <!-- 连接失败/出错: X -->
            <template v-else>
              <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
              <path d="M15 9l-6 6M9 9l6 6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            </template>
          </svg>
        </div>
      </div>
      <div class="divider"></div>
      <div class="control-group" style="cursor: pointer; --wails-draggable:no-drag" @click="$emit('quit')">
        <span class="label">❌ 退出</span>
      </div>
    </div>
  </div>

  <Teleport to="body">
    <div class="status-tooltip" v-if="showStatusTooltip" :style="tooltipStyle">
      <div class="tooltip-row">
        <span class="tooltip-label">状态:</span>
        <span class="tooltip-value">{{ statusText }}</span>
      </div>
      <div class="tooltip-row">
        <span class="tooltip-label">API状态:</span>
        <span class="tooltip-value">
          {{ statusText === '已连接' ? '✅ 接口通畅' : (statusText === 'Key无效' ? '🚫 Key无效' : (statusText === '连接失败' ? '❌ 连接失败'
            : '未配置')) }} </span>
      </div>
      <div class="tooltip-row">
        <span class="tooltip-label">模型:</span>
        <span class="tooltip-value">{{ settings.model }}</span>
      </div>
      <div class="tooltip-row">
        <span class="tooltip-label">隐身:</span>
        <span class="tooltip-value" :style="{ color: isStealthMode ? '#52c41a' : '#ff4d4f' }">
          {{ isStealthMode ? '已开启' : '已关闭' }}
        </span>
      </div>
    </div>
    <div class="settings-tooltip" v-if="showSettingsTip" :style="settingsTooltipStyle">
      <div class="tooltip-warning">
        ⚠️ 注意：打开设置将获取焦点<br>录屏期间请勿操作
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'

const props = defineProps({
  shortcuts: Object,
  activeButtons: Object,
  isClickThrough: Boolean,
  statusIcon: String,
  statusText: String,

  settings: Object,
  isStealthMode: Boolean,
  isMacOS: Boolean
})

defineEmits(['openSettings', 'quit'])

// 根据状态文本计算状态类名
const statusClass = computed(() => {
  const text = props.statusText || ''
  if (text === '已连接' || text === '就绪' || text === '解题完成') return 'connected'
  if (text.includes('未配置')) return 'unconfigured'
  if (text.includes('无效') || text.includes('Key')) return 'invalid-key'
  if (text.includes('失败') || text.includes('出错')) return 'disconnected'
  if (text.includes('思考') || text.includes('复制')) return 'connected'
  return 'unconfigured'
})

// 判断状态是否为已连接类
const isConnectedStatus = computed(() => {
  const text = props.statusText || ''
  return text === '已连接' || text === '就绪' || text === '解题完成' || text.includes('思考') || text.includes('复制')
})

// 判断是否未配置
const isUnconfigured = computed(() => {
  const text = props.statusText || ''
  return text.includes('未配置')
})

// 判断是否Key无效
const isInvalidKey = computed(() => {
  const text = props.statusText || ''
  return text.includes('无效')
})

const showStatusTooltip = ref(false)
const statusGroupRef = ref(null)
const tooltipStyle = reactive({ top: '0px', left: '0px' })

const showSettingsTip = ref(false)
const settingsBtnRef = ref(null)
const settingsTooltipStyle = reactive({ top: '0px', left: '0px' })

function showTooltip() {
  if (statusGroupRef.value) {
    const rect = statusGroupRef.value.getBoundingClientRect()
    tooltipStyle.top = `${rect.bottom + 10}px`
    tooltipStyle.left = `${rect.left + rect.width / 2}px`
    showStatusTooltip.value = true
  }
}

function hideTooltip() {
  showStatusTooltip.value = false
}

function showSettingsTooltip() {
  if (settingsBtnRef.value) {
    const rect = settingsBtnRef.value.getBoundingClientRect()
    settingsTooltipStyle.top = `${rect.bottom + 10}px`
    settingsTooltipStyle.left = `${rect.left + rect.width / 2}px`
    showSettingsTip.value = true
  }
}

function hideSettingsTooltip() {
  showSettingsTip.value = false
}
</script>

<style scoped>
/* ========================================
   TopBar Styles
   ======================================== */

.top-bar-wrapper {
  pointer-events: auto;
}

.status-group {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  padding: 0 var(--space-2);
}

.status-indicator {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--transition-fast);
}

.status-indicator.connected {
  color: var(--color-success);
}

.status-indicator.unconfigured {
  color: var(--color-warning);
}

.status-indicator.invalid-key {
  color: var(--color-error);
}

.status-indicator.disconnected {
  color: var(--text-tertiary);
}

.status-svg {
  width: 18px;
  height: 18px;
}

/* ========================================
   Tooltips
   ======================================== */

.status-tooltip {
  position: fixed;
  transform: translateX(-50%);
  background: var(--bg-elevated);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  padding: var(--space-3) var(--space-4);
  min-width: 180px;
  z-index: 99999;
  box-shadow: var(--shadow-lg);
  backdrop-filter: blur(16px);
  pointer-events: none;
  animation: tooltipIn 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.settings-tooltip {
  position: fixed;
  transform: translateX(-50%);
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.15) 0%, var(--bg-elevated) 100%);
  border: 1px solid rgba(245, 158, 11, 0.4);
  border-radius: var(--radius-md);
  padding: var(--space-3) var(--space-4);
  z-index: 99999;
  box-shadow: var(--shadow-lg), 0 0 20px rgba(245, 158, 11, 0.1);
  backdrop-filter: blur(16px);
  pointer-events: none;
  animation: tooltipIn 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  text-align: center;
}

.settings-tooltip::before {
  content: '';
  position: absolute;
  top: -6px;
  left: 50%;
  transform: translateX(-50%);
  border-width: 0 6px 6px 6px;
  border-style: solid;
  border-color: transparent transparent rgba(245, 158, 11, 0.4) transparent;
}

.tooltip-warning {
  color: var(--color-warning);
  font-size: var(--text-sm);
  line-height: 1.6;
  font-weight: 600;
}

.status-tooltip::before {
  content: '';
  position: absolute;
  top: -6px;
  left: 50%;
  transform: translateX(-50%);
  border-width: 0 6px 6px 6px;
  border-style: solid;
  border-color: transparent transparent var(--bg-elevated) transparent;
}

.tooltip-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-2);
  font-size: var(--text-sm);
  white-space: nowrap;
}

.tooltip-row:last-child {
  margin-bottom: 0;
}

.tooltip-label {
  color: var(--text-muted);
  margin-right: var(--space-4);
}

.tooltip-value {
  color: var(--text-primary);
  font-weight: 600;
  font-family: var(--font-mono);
}

@keyframes tooltipIn {
  from {
    opacity: 0;
    transform: translate(-50%, -6px);
  }
  to {
    opacity: 1;
    transform: translate(-50%, 0);
  }
}
</style>
