<template>
    <div class="model-select-container">
        <div class="model-select" :class="{ open: isOpen, disabled: disabled }" @click="toggle" ref="selectRef">

            <!-- 选中项 -->
            <div class="selected-item">
                <template v-if="modelValue">
                    <div class="provider-logo">
                        <div v-if="getProviderLogo(modelValue) && getProviderLogo(modelValue).trim().startsWith('<svg')" v-html="getProviderLogo(modelValue)"></div>
                        <img v-else-if="getProviderLogo(modelValue)" :src="getProviderLogo(modelValue)" alt="logo" />
                    </div>
                    <div class="model-info">
                        <span class="model-name">{{ modelValue }}</span>
                        <span class="provider-name">{{ getProviderName(modelValue) }}</span>
                    </div>
                    <div class="capability-tags">
                        <span v-if="getModelCapabilities(modelValue).image" class="cap-tag" title="支持图片">
                            <svg class="cap-icon" viewBox="0 0 16 16" fill="none">
                                <rect x="2" y="3" width="12" height="10" rx="1.5" stroke="currentColor" stroke-width="1.2"/>
                                <circle cx="5.5" cy="6.5" r="1" fill="currentColor"/>
                                <path d="M2.5 11l3-3 2 2 4-4 2 2" stroke="currentColor" stroke-width="1.2" stroke-linecap="round" stroke-linejoin="round"/>
                            </svg>
                        </span>
                        <span v-if="getModelCapabilities(modelValue).file" class="cap-tag" title="支持文件">
                            <svg class="cap-icon" viewBox="0 0 16 16" fill="none">
                                <path d="M4 2h5l4 4v8a1 1 0 01-1 1H4a1 1 0 01-1-1V3a1 1 0 011-1z" stroke="currentColor" stroke-width="1.2"/>
                                <path d="M9 2v4h4" stroke="currentColor" stroke-width="1.2"/>
                            </svg>
                        </span>
                        <span v-if="getModelCapabilities(modelValue).audio" class="cap-tag" title="支持音频">
                            <svg class="cap-icon" viewBox="0 0 16 16" fill="none">
                                <path d="M8 3v10M5 5.5v5M11 5.5v5M2 7v2M14 7v2" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                            </svg>
                        </span>
                        <span v-if="getModelCapabilities(modelValue).video" class="cap-tag" title="支持视频">
                            <svg class="cap-icon" viewBox="0 0 16 16" fill="none">
                                <rect x="1" y="4" width="10" height="8" rx="1" stroke="currentColor" stroke-width="1.2"/>
                                <path d="M11 7l4-2v6l-4-2v-2z" stroke="currentColor" stroke-width="1.2" stroke-linejoin="round"/>
                            </svg>
                        </span>
                        <span v-if="getModelCapabilities(modelValue).contextLength" class="cap-tag cap-context" title="上下文长度">
                            {{ formatContextLength(getModelCapabilities(modelValue).contextLength) }}
                        </span>
                    </div>
                </template>
                <span v-else class="placeholder">请选择模型</span>
                <span class="arrow" :class="{ rotated: isOpen }">
                    <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
                        <path d="M2.5 4.5L6 8L9.5 4.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"
                            stroke-linejoin="round" />
                    </svg>
                </span>
            </div>

            <!-- 下拉列表 -->
            <Transition name="dropdown">
                <div v-if="isOpen" class="dropdown-list">
                    <div v-if="loading" class="loading-state">
                        <span class="loading-icon">⏳</span>
                        <span>加载中...</span>
                    </div>
                    <template v-else>
                        <div v-for="model in models" :key="model" class="dropdown-item"
                            :class="{ selected: modelValue === model }" @click.stop="selectModel(model)">
                            <div class="provider-logo">
                                <div v-if="getProviderLogo(model) && getProviderLogo(model).trim().startsWith('<svg')" v-html="getProviderLogo(model)"></div>
                                <img v-else-if="getProviderLogo(model)" :src="getProviderLogo(model)" alt="logo" />
                            </div>
                            <div class="model-info">
                                <span class="model-name">{{ model }}</span>
                                <span class="provider-name">{{ getProviderName(model) }}</span>
                            </div>
                            <div class="capability-tags">
                                <span v-if="getModelCapabilities(model).image" class="cap-tag" title="支持图片">
                                    <svg class="cap-icon" viewBox="0 0 16 16" fill="none">
                                        <rect x="2" y="3" width="12" height="10" rx="1.5" stroke="currentColor" stroke-width="1.2"/>
                                        <circle cx="5.5" cy="6.5" r="1" fill="currentColor"/>
                                        <path d="M2.5 11l3-3 2 2 4-4 2 2" stroke="currentColor" stroke-width="1.2" stroke-linecap="round" stroke-linejoin="round"/>
                                    </svg>
                                </span>
                                <span v-if="getModelCapabilities(model).file" class="cap-tag" title="支持文件">
                                    <svg class="cap-icon" viewBox="0 0 16 16" fill="none">
                                        <path d="M4 2h5l4 4v8a1 1 0 01-1 1H4a1 1 0 01-1-1V3a1 1 0 011-1z" stroke="currentColor" stroke-width="1.2"/>
                                        <path d="M9 2v4h4" stroke="currentColor" stroke-width="1.2"/>
                                    </svg>
                                </span>
                                <span v-if="getModelCapabilities(model).audio" class="cap-tag" title="支持音频">
                                    <svg class="cap-icon" viewBox="0 0 16 16" fill="none">
                                        <path d="M8 3v10M5 5.5v5M11 5.5v5M2 7v2M14 7v2" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                                    </svg>
                                </span>
                                <span v-if="getModelCapabilities(model).video" class="cap-tag" title="支持视频">
                                    <svg class="cap-icon" viewBox="0 0 16 16" fill="none">
                                        <rect x="1" y="4" width="10" height="8" rx="1" stroke="currentColor" stroke-width="1.2"/>
                                        <path d="M11 7l4-2v6l-4-2v-2z" stroke="currentColor" stroke-width="1.2" stroke-linejoin="round"/>
                                    </svg>
                                </span>
                                <span v-if="getModelCapabilities(model).contextLength" class="cap-tag cap-context" title="上下文长度">
                                    {{ formatContextLength(getModelCapabilities(model).contextLength) }}
                                </span>
                            </div>
                            <span v-if="modelValue === model" class="check-icon">✓</span>
                        </div>
                    </template>
                </div>
            </Transition>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { getProviderLogo, getProviderName, getModelCapabilities, formatContextLength } from '../utils/modelCapabilities'

const props = defineProps({
    modelValue: {
        type: String,
        default: ''
    },
    models: {
        type: Array,
        default: () => []
    },
    loading: {
        type: Boolean,
        default: false
    },
    disabled: {
        type: Boolean,
        default: false
    }
})

const emit = defineEmits(['update:modelValue'])

const isOpen = ref(false)
const selectRef = ref(null)

function toggle() {
    if (props.disabled || props.loading) return
    isOpen.value = !isOpen.value
}

function selectModel(model) {
    emit('update:modelValue', model)
    isOpen.value = false
}

// 点击外部关闭
function handleClickOutside(event) {
    if (selectRef.value && !selectRef.value.contains(event.target)) {
        isOpen.value = false
    }
}

onMounted(() => {
    document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
    document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.model-select-container {
    width: 100%;
}

.model-select {
    position: relative;
    background: var(--bg-elevated);
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-md);
    cursor: pointer;
    transition: all var(--transition-fast);
}

.model-select:hover {
    border-color: var(--border-default);
    background: var(--bg-card-hover);
}

.model-select.open {
    border-color: var(--color-primary);
    box-shadow: var(--shadow-focus);
}

.model-select.disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.selected-item {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    padding: var(--space-3) var(--space-4);
}

.provider-logo {
    width: 28px;
    height: 28px;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--text-primary);
    border-radius: var(--radius-sm);
    overflow: hidden;
}

.provider-logo :deep(svg) {
    width: 100%;
    height: 100%;
    display: block;
    max-width: 28px;
    max-height: 28px;
}

.provider-logo img {
    width: 100%;
    height: 100%;
    object-fit: contain;
}

.model-info {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.model-name {
    font-size: var(--text-sm);
    font-weight: 500;
    color: var(--text-primary);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.provider-name {
    font-size: var(--text-xs);
    color: var(--text-tertiary);
}

/* 能力标签 - 精致图标风格 */
.capability-tags {
    display: flex;
    gap: 4px;
    flex-shrink: 0;
    align-items: center;
}

.cap-tag {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 22px;
    height: 22px;
    border-radius: 5px;
    background: rgba(99, 102, 241, 0.1);
    color: var(--color-primary);
    transition: all var(--transition-fast);
}

.cap-tag:hover {
    background: rgba(99, 102, 241, 0.2);
}

.cap-tag.cap-context {
    width: auto;
    padding: 0 6px;
    font-size: 10px;
    font-weight: 600;
    font-family: var(--font-mono, 'SF Mono', 'Menlo', monospace);
    letter-spacing: -0.02em;
}

.cap-icon {
    width: 13px;
    height: 13px;
}

.placeholder {
    color: var(--text-tertiary);
    font-size: var(--text-sm);
}

.arrow {
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--text-tertiary);
    transition: transform var(--transition-fast);
    margin-left: auto;
}

.arrow.rotated {
    transform: rotate(180deg);
}

/* 下拉列表 */
.dropdown-list {
    position: absolute;
    top: calc(100% + 4px);
    left: 0;
    right: 0;
    max-height: 320px;
    overflow-y: auto;
    background: var(--bg-elevated);
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-xl);
    z-index: 100;
    backdrop-filter: blur(20px);
}

.dropdown-list::-webkit-scrollbar {
    width: 6px;
}

.dropdown-list::-webkit-scrollbar-track {
    background: transparent;
}

.dropdown-list::-webkit-scrollbar-thumb {
    background: var(--scrollbar-thumb);
    border-radius: 3px;
}

.dropdown-item {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    padding: var(--space-3) var(--space-4);
    cursor: pointer;
    transition: background var(--transition-fast);
    border-bottom: 1px solid var(--border-subtle);
}

.dropdown-item:last-child {
    border-bottom: none;
}

.dropdown-item:hover {
    background: var(--bg-card-hover);
}

.dropdown-item.selected {
    background: var(--color-primary-light);
}

.dropdown-item .provider-logo {
    width: 24px;
    height: 24px;
}

.dropdown-item .model-name {
    font-size: var(--text-sm);
}

.dropdown-item .provider-name {
    font-size: 10px;
}

.dropdown-item .cap-tag {
    width: 18px;
    height: 18px;
    border-radius: 4px;
}

.dropdown-item .cap-tag.cap-context {
    width: auto;
    padding: 0 5px;
    font-size: 9px;
}

.dropdown-item .cap-icon {
    width: 11px;
    height: 11px;
}

.check-icon {
    color: var(--color-primary);
    font-weight: bold;
    margin-left: auto;
}

.loading-state {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: var(--space-2);
    padding: var(--space-5);
    color: var(--text-tertiary);
}

.loading-icon {
    animation: selectSpin 1s linear infinite;
}

@keyframes selectSpin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
}

/* 下拉动画 */
.dropdown-enter-active,
.dropdown-leave-active {
    transition: all var(--transition-fast);
}

.dropdown-enter-from,
.dropdown-leave-to {
    opacity: 0;
    transform: translateY(-8px);
}
</style>
