<template>
    <div class="model-select-container">
        <div class="model-select" :class="{ open: isOpen, disabled: disabled }" @click="toggle" ref="selectRef">

            <!-- 选中项 -->
            <div class="selected-item">
                <template v-if="modelValue">
                    <div class="provider-logo" v-html="getProviderLogo(modelValue)"></div>
                    <div class="model-info">
                        <span class="model-name">{{ modelValue }}</span>
                        <span class="provider-name">{{ getProviderName(modelValue) }}</span>
                    </div>
                    <div class="capability-badges">
                        <span v-if="getModelCapabilities(modelValue).image" class="cap-badge" title="支持图片">🖼️</span>
                        <span v-if="getModelCapabilities(modelValue).pdf" class="cap-badge" title="支持PDF">📄</span>
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
                            <div class="provider-logo" v-html="getProviderLogo(model)"></div>
                            <div class="model-info">
                                <span class="model-name">{{ model }}</span>
                                <span class="provider-name">{{ getProviderName(model) }}</span>
                            </div>
                            <div class="capability-badges">
                                <span v-if="getModelCapabilities(model).image" class="cap-badge" title="支持图片">🖼️</span>
                                <span v-if="getModelCapabilities(model).pdf" class="cap-badge" title="支持PDF">📄</span>
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
import { getProviderLogo, getProviderName, getModelCapabilities } from '../utils/modelCapabilities'

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
    background: rgba(30, 30, 36, 0.95);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 12px;
    cursor: pointer;
    transition: all 0.2s ease;
}

.model-select:hover {
    border-color: rgba(255, 255, 255, 0.2);
    background: rgba(40, 40, 48, 0.95);
}

.model-select.open {
    border-color: #4CAF50;
    box-shadow: 0 0 0 2px rgba(76, 175, 80, 0.2);
}

.model-select.disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.selected-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
}

.provider-logo {
    width: 28px;
    height: 28px;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
}

.provider-logo :deep(svg) {
    width: 100%;
    height: 100%;
}

.model-info {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.model-name {
    font-size: 14px;
    font-weight: 500;
    color: #fff;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.provider-name {
    font-size: 11px;
    color: rgba(255, 255, 255, 0.5);
}

.capability-badges {
    display: flex;
    gap: 4px;
    flex-shrink: 0;
}

.cap-badge {
    font-size: 14px;
    opacity: 0.8;
}

.placeholder {
    color: rgba(255, 255, 255, 0.4);
    font-size: 14px;
}

.arrow {
    display: flex;
    align-items: center;
    justify-content: center;
    color: rgba(255, 255, 255, 0.5);
    transition: transform 0.2s ease;
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
    background: rgba(28, 28, 34, 0.98);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 12px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
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
    background: rgba(255, 255, 255, 0.2);
    border-radius: 3px;
}

.dropdown-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px 16px;
    cursor: pointer;
    transition: background 0.15s ease;
    border-bottom: 1px solid rgba(255, 255, 255, 0.03);
}

.dropdown-item:last-child {
    border-bottom: none;
}

.dropdown-item:hover {
    background: rgba(255, 255, 255, 0.08);
}

.dropdown-item.selected {
    background: rgba(76, 175, 80, 0.15);
}

.dropdown-item .provider-logo {
    width: 24px;
    height: 24px;
}

.dropdown-item .model-name {
    font-size: 13px;
}

.dropdown-item .provider-name {
    font-size: 10px;
}

.check-icon {
    color: #4CAF50;
    font-weight: bold;
    margin-left: auto;
}

.loading-state {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    padding: 20px;
    color: rgba(255, 255, 255, 0.5);
}

.loading-icon {
    animation: spin 1s linear infinite;
}

@keyframes spin {
    from {
        transform: rotate(0deg);
    }

    to {
        transform: rotate(360deg);
    }
}

/* 下拉动画 */
.dropdown-enter-active,
.dropdown-leave-active {
    transition: all 0.2s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
    opacity: 0;
    transform: translateY(-8px);
}
</style>
