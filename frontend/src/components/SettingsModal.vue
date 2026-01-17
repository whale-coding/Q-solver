<template>
  <div v-if="show" class="modal" id="settings-modal" style="display: flex">
    <div class="modal-content">
      <div class="modal-warning-banner"
        style="background: rgba(255, 169, 64, 0.15); border: 1px solid rgba(255, 169, 64, 0.3); border-radius: 50px; padding: 6px 20px; color: #ffc069; font-size: 12px; display: flex; align-items: center; justify-content: center; margin: 12px auto 4px auto; width: fit-content;">
        ⚠️ 当前窗口已获取焦点，关闭设置后将自动恢复防抢焦模式
      </div>
      <div class="modal-header">
        <div class="tabs">
          <div class="tab" :class="{ active: currentTab === 'general' }" @click="currentTab = 'general'">
            常规设置</div>
          <div class="tab" :class="{ active: currentTab === 'model' }" @click="currentTab = 'model'">模型设置
          </div>
          <div class="tab" :class="{ active: currentTab === 'params' }" @click="currentTab = 'params'">生成参数</div>
          <div class="tab" :class="{ active: currentTab === 'screenshot' }" @click="currentTab = 'screenshot'">截图设置</div>
          <div class="tab" :class="{ active: currentTab === 'resume' }" @click="currentTab = 'resume'">
            简历设置</div>
          <div class="tab" :class="{ active: currentTab === 'account' }" @click="currentTab = 'account'">
            提供商</div>
        </div>
        <span class="close-btn" @click="$emit('close')">&times;</span>
      </div>
      <div class="modal-body">
        <div v-show="currentTab === 'account'">
          <ProviderSelect v-model:provider="tempSettings.provider" v-model:apiKey="tempSettings.apiKey"
            v-model:baseURL="tempSettings.baseURL" />
        </div>

        <div v-show="currentTab === 'model'">
          <div class="form-group">
            <div class="model-header">
              <label>模型选择</label>
              <div class="model-actions">
                <button class="btn-icon" @click="$emit('refresh-models')"
                  :disabled="isLoadingModels || !tempSettings.apiKey" title="刷新模型列表">
                  <svg class="action-icon" :class="{ spin: isLoadingModels }" viewBox="0 0 16 16" fill="none">
                    <path d="M14 8a6 6 0 01-10.24 4.24" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                    <path d="M2 8a6 6 0 0110.24-4.24" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                    <path d="M14 3v5h-5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                    <path d="M2 13V8h5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                </button>
                <button class="btn-icon" @click="$emit('test-connection')"
                  :disabled="isTestingConnection || !tempSettings.model" title="测试模型连通性">
                  <svg v-if="isTestingConnection" class="action-icon spin" viewBox="0 0 16 16" fill="none">
                    <circle cx="8" cy="8" r="6" stroke="currentColor" stroke-width="1.5" stroke-dasharray="28 10" stroke-linecap="round"/>
                  </svg>
                  <svg v-else class="action-icon" viewBox="0 0 16 16" fill="none">
                    <path d="M4 3l9 5-9 5V3z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/>
                  </svg>
                </button>
              </div>
            </div>
            <ModelSelect v-model="tempSettings.model" :models="availableModels" :loading="isLoadingModels" />

            <!-- 连通性测试结果 -->
            <div v-if="connectionStatus" class="connection-status" :class="connectionStatus.type">
              <span class="status-icon">{{ connectionStatus.icon }}</span>
              <span class="status-text">{{ connectionStatus.message }}</span>
            </div>

            <p v-if="!tempSettings.apiKey" class="hint-text warning-hint">
              ⚠️ 请先填写 API Key
            </p>
          </div>

          <div class="form-group">
            <div class="prompt-header">
              <label for="prompt-text" style="margin-bottom: 0">系统提示词 (Prompt)</label>
              <div class="prompt-tabs">
                <div class="prompt-tab" :class="{ active: promptTab === 'edit' }" @click="promptTab = 'edit'">编辑
                </div>
                <div class="prompt-tab" :class="{ active: promptTab === 'preview' }" @click="promptTab = 'preview'">预览
                </div>
              </div>
            </div>

            <textarea v-show="promptTab === 'edit'" id="prompt-text" class="prompt-textarea" rows="10"
              v-model="tempSettings.prompt" placeholder="请输入提示词 (支持 Markdown)..."></textarea>

            <div v-show="promptTab === 'preview'" class="prompt-preview markdown-body" v-html="renderedPrompt">
            </div>
          </div>
        </div>

        <div v-show="currentTab === 'params'">
          <LLMParamsConfig v-model:temperature="tempSettings.temperature" v-model:topP="tempSettings.topP"
            v-model:topK="tempSettings.topK" v-model:maxTokens="tempSettings.maxTokens"
            v-model:thinkingBudget="tempSettings.thinkingBudget" />
        </div>

        <div v-show="currentTab === 'general'">
          <div class="form-group">
            <div class="context-setting">
              <div class="setting-row">
                <div class="setting-info">
                  <span class="setting-title">保存上下文</span>
                  <span class="setting-desc">开启后，每次对话将包含之前的历史记录</span>
                </div>
                <label class="switch">
                  <input type="checkbox" v-model="tempSettings.keepContext">
                  <span class="slider round"></span>
                </label>
              </div>

              <div class="setting-row" style="margin-top: 12px;">
                <div class="setting-info">
                  <span class="setting-title">启用 Live API 模式</span>
                  <span class="setting-desc">采集扬声器声音，实时识别面试官问题并回答</span>
                </div>
                <label class="switch">
                  <input type="checkbox" v-model="tempSettings.useLiveApi">
                  <span class="slider round"></span>
                </label>
              </div>
            </div>
          </div>

          <div class="form-group">
            <label>快捷键配置 {{ isMacOS ? '(macOS 不支持自定义)' : '(点击录制)' }}</label>
            <div class="shortcut-list">
              <div class="shortcut-item" v-for="key in shortcutActions" :key="key.action">
                <span>{{ key.label }}</span>
                <button class="btn-record" :class="{ recording: recordingAction === key.action, disabled: isMacOS }"
                  @click="!isMacOS && $emit('record-key', key.action)"
                  :title="isMacOS ? 'macOS 使用预设快捷键，不支持自定义' : '点击录制新快捷键'">
                  {{ recordingAction === key.action ? recordingText : (tempShortcuts[key.action]?.keyName ||
                    (isMacOS ? key.macDefault : key.default)) }}
                </button>
              </div>
            </div>
          </div>

          <div class="form-group">
            <label for="opacity-slider">窗口透明度: <span>{{ Math.round(tempSettings.transparency * 100) }}%</span></label>
            <input type="range" id="opacity-slider" min="0.0" max="1.0" step="0.05"
              v-model.number="tempSettings.transparency" />
          </div>
        </div>

        <div v-show="currentTab === 'screenshot'">
          <ScreenshotSettings :modelValue="tempSettings" @update:modelValue="Object.assign(tempSettings, $event)" />
        </div>

        <div v-show="currentTab === 'resume'" style="height: 100%">
          <ResumeImport :resumePath="tempSettings.resumePath" :rawContent="resumeRawContent"
            :isParsing="isResumeParsing" :currentModel="tempSettings.model"
            v-model:useMarkdownResume="tempSettings.useMarkdownResume"
            @update:rawContent="$emit('update:resumeRawContent', $event)" @select-resume="$emit('select-resume')"
            @clear-resume="$emit('clear-resume')" @parse-resume="$emit('parse-resume')" />
        </div>
      </div>
      <div class="modal-footer">
        <button class="btn-primary" @click="$emit('save')">保存</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import ResumeImport from './ResumeImport.vue'
import ScreenshotSettings from './ScreenshotSettings.vue'
import ProviderSelect from './ProviderSelect.vue'
import ModelSelect from './ModelSelect.vue'
import LLMParamsConfig from './LLMParamsConfig.vue'

const props = defineProps({
  show: Boolean,
  tempSettings: Object,
  tempShortcuts: Object,
  shortcutActions: Array,
  recordingAction: String,
  recordingText: String,
  availableModels: Array,
  isLoadingModels: Boolean,
  isTestingConnection: Boolean,
  connectionStatus: Object,
  renderedPrompt: String,
  resumeRawContent: String,
  isResumeParsing: Boolean,
  isMacOS: Boolean,
  activeTab: {
    type: String,
    defaut: 'general'
  }
})

const emit = defineEmits([
  'close',
  'save',
  'refresh-models',
  'test-connection',
  'record-key',
  'select-resume',
  'clear-resume',
  'parse-resume',
  'update:resumeRawContent',
  'update:activeTab'
])

const currentTab = computed({
  get: () => props.activeTab || 'general',
  set: (val) => emit('update:activeTab', val)
})

const promptTab = ref('edit')
</script>
