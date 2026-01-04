<template>
    <div class="resume-import">
        <!-- Empty State - No file selected -->
        <div v-if="!resumePath" class="empty-state">
            <div class="empty-state-card" @click="$emit('select-resume')">
                <div class="upload-icon-wrapper">
                    <div class="upload-icon">📄</div>
                    <div class="upload-pulse"></div>
                </div>
                <h3 class="empty-title">导入 PDF 简历</h3>
                <p class="empty-desc">AI 将在解题时参考您的背景信息，提供更个性化的回答</p>
                <button class="btn-upload">
                    <span class="btn-icon">📂</span>
                    选择文件
                </button>
                <p class="empty-hint">支持 .pdf 格式</p>
            </div>
        </div>

        <!-- File Selected State -->
        <div v-else class="resume-content">
            <!-- File Info Header -->
            <div class="file-header">
                <div class="file-info">
                    <div class="file-icon">📎</div>
                    <div class="file-details">
                        <span class="file-name">{{ fileName }}</span>
                    </div>
                </div>
                <div class="file-actions">
                    <!-- Use Markdown Toggle -->
                    <div class="toggle-chip" :class="{ active: useMarkdownResume }"
                        @click="$emit('update:useMarkdownResume', !useMarkdownResume)" title="使用解析后的 Markdown 文本">
                        <span class="toggle-dot"></span>
                        <span>Markdown 模式</span>
                    </div>
                    <!-- Menu Button -->
                    <div class="menu-wrapper">
                        <button class="btn-menu" @click="showMenu = !showMenu">⋮</button>
                        <div v-if="showMenu" class="dropdown-menu">
                            <div class="menu-item" @click="handleMenuAction('change')">
                                <span>📂</span> 更换文件
                            </div>
                            <div class="menu-item danger" @click="handleMenuAction('clear')">
                                <span>🗑️</span> 清除简历
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Content Area -->
            <div class="content-area">
                <!-- Status Card (no markdown yet) -->
                <div v-if="!localContent && !isEditing" class="status-card">
                    <div class="status-icon">📄</div>
                    <h4 class="status-title">文件已就绪</h4>
                    <p class="status-desc">点击 AI 解析将 PDF 转换为 Markdown 格式</p>

                    <div class="status-actions">
                        <button class="btn-parse-lg" @click="handleParseClick" :disabled="isParsing">
                            <span v-if="!isParsing">✨</span>
                            <span v-else class="spin">⏳</span>
                            {{ isParsing ? '解析中...' : 'AI 解析为 Markdown' }}
                        </button>
                        <button class="btn-secondary-lg" @click="isEditing = true">
                            📝 手动输入
                        </button>
                    </div>

                    <p v-if="!modelSupportsFile" class="status-warning">
                        ⚠️ 当前模型可能不支持 PDF 解析
                    </p>
                </div>

                <!-- Parsing State -->
                <div v-else-if="isParsing && !localContent" class="status-card parsing">
                    <span class="spin large">⏳</span>
                    <h4 class="status-title">AI 正在解析</h4>
                    <p class="status-desc">正在将 PDF 转换为 Markdown...</p>
                </div>

                <!-- Markdown Content -->
                <div v-else class="markdown-panel">
                    <!-- Toolbar -->
                    <div class="markdown-toolbar">
                        <div class="toolbar-tabs">
                            <button class="toolbar-btn" :class="{ active: !isEditing }" @click="isEditing = false">
                                预览
                            </button>
                            <button class="toolbar-btn" :class="{ active: isEditing }" @click="isEditing = true">
                                编辑
                            </button>
                        </div>
                        <div class="toolbar-actions">
                            <button class="btn-reparse" @click="handleParseClick" :disabled="isParsing" title="重新解析">
                                <span v-if="!isParsing">🔄</span>
                                <span v-else class="spin">⏳</span>
                            </button>
                        </div>
                    </div>

                    <!-- Editor -->
                    <div v-if="isEditing" class="editor-wrapper">
                        <textarea v-model="localContent" @input="updateContent" class="md-editor"
                            placeholder="在此输入或粘贴 Markdown 格式的简历内容..."></textarea>
                    </div>

                    <!-- Preview -->
                    <div v-else class="md-preview" v-html="renderedContent"></div>
                </div>
            </div>
        </div>

        <!-- Confirm Dialog -->
        <Teleport to="body">
            <div v-if="showConfirmDialog" class="dialog-overlay" @click.self="showConfirmDialog = false">
                <div class="dialog-box">
                    <div class="dialog-icon">⚠️</div>
                    <h4>模型可能不支持</h4>
                    <p>当前模型可能不支持 PDF 解析，是否仍要继续？</p>
                    <div class="dialog-actions">
                        <button class="btn-cancel" @click="showConfirmDialog = false">取消</button>
                        <button class="btn-confirm" @click="confirmParse">继续</button>
                    </div>
                </div>
            </div>
        </Teleport>
    </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue';
import { marked } from 'marked';
import { supportsVision, supportsPDF } from '../utils/modelCapabilities';

const props = defineProps({
    resumePath: { type: String, default: '' },
    rawContent: { type: String, default: '' },
    isParsing: { type: Boolean, default: false },
    useMarkdownResume: { type: Boolean, default: false },
    currentModel: { type: String, default: '' }
});

const emit = defineEmits(['select-resume', 'clear-resume', 'parse-resume', 'update:rawContent', 'update:useMarkdownResume']);

// UI State
const isEditing = ref(false);
const showMenu = ref(false);
const showConfirmDialog = ref(false);
const localContent = ref(props.rawContent);

// Computed
const modelSupportsFile = computed(() => supportsVision(props.currentModel) || supportsPDF(props.currentModel));
const fileName = computed(() => {
    if (!props.resumePath) return '';
    return props.resumePath.split(/[\\/]/).pop() || 'resume.pdf';
});
const renderedContent = computed(() => {
    if (!localContent.value) return '';
    return marked.parse(localContent.value);
});

// Watchers
watch(() => props.rawContent, (newVal) => {
    if (newVal !== localContent.value) {
        localContent.value = newVal;
    }
});

// Click outside to close menu
watch(showMenu, (val) => {
    if (val) {
        setTimeout(() => {
            document.addEventListener('click', closeMenu);
        }, 0);
    }
});

function closeMenu() {
    showMenu.value = false;
    document.removeEventListener('click', closeMenu);
}

// Content Functions
function updateContent() {
    emit('update:rawContent', localContent.value);
}

// Action Handlers
function handleParseClick() {
    if (!modelSupportsFile.value) {
        showConfirmDialog.value = true;
    } else {
        emit('parse-resume');
    }
}

function confirmParse() {
    showConfirmDialog.value = false;
    emit('parse-resume');
}

function handleMenuAction(action) {
    showMenu.value = false;
    switch (action) {
        case 'change':
            emit('select-resume');
            break;
        case 'clear':
            emit('clear-resume');
            break;
    }
}
</script>

<style scoped>
/* ========================================
   Resume Import - Simplified UI
   ======================================== */

.resume-import {
    height: 100%;
    display: flex;
    flex-direction: column;
    color: #fff;
}

/* ========================================
   Empty State
   ======================================== */

.empty-state {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 20px;
}

.empty-state-card {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 40px 50px;
    background: linear-gradient(145deg, rgba(40, 40, 50, 0.6), rgba(30, 30, 40, 0.8));
    border: 2px dashed rgba(99, 102, 241, 0.3);
    border-radius: 16px;
    cursor: pointer;
    transition: all 0.3s ease;
}

.empty-state-card:hover {
    border-color: rgba(99, 102, 241, 0.6);
    background: linear-gradient(145deg, rgba(50, 50, 65, 0.7), rgba(35, 35, 50, 0.9));
    transform: translateY(-2px);
}

.upload-icon-wrapper {
    position: relative;
    margin-bottom: 20px;
}

.upload-icon {
    font-size: 48px;
    position: relative;
    z-index: 1;
}

.upload-pulse {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 80px;
    height: 80px;
    background: rgba(99, 102, 241, 0.15);
    border-radius: 50%;
    animation: pulse 2s infinite;
}

@keyframes pulse {

    0%,
    100% {
        transform: translate(-50%, -50%) scale(1);
        opacity: 0.5;
    }

    50% {
        transform: translate(-50%, -50%) scale(1.15);
        opacity: 0.2;
    }
}

.empty-title {
    font-size: 18px;
    font-weight: 600;
    margin: 0 0 8px 0;
    color: #fff;
}

.empty-desc {
    font-size: 13px;
    color: rgba(255, 255, 255, 0.5);
    margin: 0 0 20px 0;
    text-align: center;
    max-width: 240px;
}

.btn-upload {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 24px;
    background: linear-gradient(135deg, #6366f1, #8b5cf6);
    border: none;
    border-radius: 8px;
    color: #fff;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
}

.btn-upload:hover {
    transform: scale(1.02);
    box-shadow: 0 4px 20px rgba(99, 102, 241, 0.4);
}

.empty-hint {
    font-size: 11px;
    color: rgba(255, 255, 255, 0.3);
    margin: 12px 0 0 0;
}

/* ========================================
   File Header
   ======================================== */

.resume-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 12px;
    min-height: 0;
}

.file-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    background: rgba(255, 255, 255, 0.03);
    border-radius: 10px;
    border: 1px solid rgba(255, 255, 255, 0.06);
}

.file-info {
    display: flex;
    align-items: center;
    gap: 12px;
}

.file-icon {
    font-size: 20px;
    width: 36px;
    height: 36px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(99, 102, 241, 0.1);
    border-radius: 8px;
}

.file-details {
    display: flex;
    align-items: center;
}

.file-name {
    font-size: 13px;
    font-weight: 500;
    color: #fff;
}

.file-actions {
    display: flex;
    align-items: center;
    gap: 8px;
}

.toggle-chip {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 6px 12px;
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 20px;
    font-size: 11px;
    color: rgba(255, 255, 255, 0.5);
    cursor: pointer;
    transition: all 0.2s;
}

.toggle-chip:hover {
    background: rgba(255, 255, 255, 0.08);
}

.toggle-chip.active {
    background: rgba(99, 102, 241, 0.15);
    border-color: rgba(99, 102, 241, 0.3);
    color: #a5b4fc;
}

.toggle-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.3);
    transition: all 0.2s;
}

.toggle-chip.active .toggle-dot {
    background: #6366f1;
    box-shadow: 0 0 6px rgba(99, 102, 241, 0.5);
}

.menu-wrapper {
    position: relative;
}

.btn-menu {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 6px;
    color: rgba(255, 255, 255, 0.6);
    font-size: 16px;
    cursor: pointer;
    transition: all 0.2s;
}

.btn-menu:hover {
    background: rgba(255, 255, 255, 0.1);
}

.dropdown-menu {
    position: absolute;
    top: 100%;
    right: 0;
    margin-top: 4px;
    min-width: 140px;
    background: #2a2a35;
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 8px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
    z-index: 100;
    overflow: hidden;
}

.menu-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 14px;
    font-size: 12px;
    color: rgba(255, 255, 255, 0.7);
    cursor: pointer;
    transition: background 0.15s;
}

.menu-item:hover {
    background: rgba(255, 255, 255, 0.05);
}

.menu-item.danger {
    color: #f87171;
}

.menu-item.danger:hover {
    background: rgba(248, 113, 113, 0.1);
}

/* ========================================
   Content Area
   ======================================== */

.content-area {
    flex: 1;
    min-height: 0;
    border-radius: 10px;
    overflow: hidden;
    background: rgba(0, 0, 0, 0.2);
    border: 1px solid rgba(255, 255, 255, 0.05);
}

/* ========================================
   Status Card
   ======================================== */

.status-card {
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 12px;
    padding: 40px;
    text-align: center;
}

.status-card.parsing {
    gap: 16px;
}

.status-icon {
    font-size: 48px;
    opacity: 0.8;
}

.status-title {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
    color: #fff;
}

.status-desc {
    margin: 0;
    font-size: 13px;
    color: rgba(255, 255, 255, 0.5);
}

.status-actions {
    display: flex;
    gap: 16px;
    margin-top: 16px;
}

.btn-parse-lg {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    padding: 10px 20px;
    background: transparent;
    border: 1.5px solid rgba(99, 102, 241, 0.5);
    border-radius: 8px;
    color: #a5b4fc;
    font-size: 13px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
}

.btn-parse-lg:hover:not(:disabled) {
    background: rgba(99, 102, 241, 0.1);
    border-color: rgba(99, 102, 241, 0.8);
    color: #c7d2fe;
}

.btn-parse-lg:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.btn-secondary-lg {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    padding: 10px 20px;
    background: transparent;
    border: 1.5px solid rgba(255, 255, 255, 0.15);
    border-radius: 8px;
    color: rgba(255, 255, 255, 0.6);
    font-size: 13px;
    cursor: pointer;
    transition: all 0.2s ease;
}

.btn-secondary-lg:hover {
    background: rgba(255, 255, 255, 0.05);
    border-color: rgba(255, 255, 255, 0.25);
    color: rgba(255, 255, 255, 0.8);
}

.status-warning {
    margin: 8px 0 0 0;
    font-size: 12px;
    color: #fbbf24;
}

/* ========================================
   Markdown Panel
   ======================================== */

.markdown-panel {
    height: 100%;
    display: flex;
    flex-direction: column;
}

.markdown-toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 12px;
    background: rgba(0, 0, 0, 0.2);
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.toolbar-tabs {
    display: flex;
    gap: 2px;
}

.toolbar-btn {
    padding: 6px 14px;
    background: transparent;
    border: none;
    border-radius: 4px;
    color: rgba(255, 255, 255, 0.5);
    font-size: 12px;
    cursor: pointer;
    transition: all 0.15s;
}

.toolbar-btn:hover {
    color: rgba(255, 255, 255, 0.7);
}

.toolbar-btn.active {
    background: rgba(99, 102, 241, 0.15);
    color: #a5b4fc;
}

.btn-reparse {
    width: 28px;
    height: 28px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(255, 255, 255, 0.05);
    border: none;
    border-radius: 4px;
    color: rgba(255, 255, 255, 0.5);
    font-size: 14px;
    cursor: pointer;
    transition: all 0.2s;
}

.btn-reparse:hover:not(:disabled) {
    background: rgba(255, 255, 255, 0.1);
    color: rgba(255, 255, 255, 0.8);
}

.btn-reparse:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.editor-wrapper {
    flex: 1;
    min-height: 0;
}

.md-editor {
    width: 100%;
    height: 100%;
    padding: 16px;
    background: transparent;
    border: none;
    color: #e0e0e0;
    font-family: 'Fira Code', 'Consolas', monospace;
    font-size: 13px;
    line-height: 1.6;
    resize: none;
    outline: none;
}

.md-editor::placeholder {
    color: rgba(255, 255, 255, 0.25);
}

.md-preview {
    flex: 1;
    padding: 16px;
    overflow-y: auto;
    font-size: 13px;
    line-height: 1.6;
    color: #e0e0e0;
}

.md-preview :deep(h1),
.md-preview :deep(h2),
.md-preview :deep(h3) {
    margin: 1em 0 0.5em 0;
    color: #fff;
}

.md-preview :deep(h1) {
    font-size: 1.4em;
}

.md-preview :deep(h2) {
    font-size: 1.2em;
}

.md-preview :deep(h3) {
    font-size: 1.1em;
}

.md-preview :deep(p) {
    margin: 0.5em 0;
}

.md-preview :deep(ul),
.md-preview :deep(ol) {
    padding-left: 1.5em;
    margin: 0.5em 0;
}

.md-preview :deep(code) {
    background: rgba(255, 255, 255, 0.1);
    padding: 2px 6px;
    border-radius: 4px;
    font-size: 0.9em;
}

.md-preview :deep(pre) {
    background: rgba(0, 0, 0, 0.3);
    padding: 12px;
    border-radius: 6px;
    overflow-x: auto;
}

.md-preview :deep(blockquote) {
    border-left: 3px solid #6366f1;
    padding-left: 12px;
    margin: 0.5em 0;
    color: rgba(255, 255, 255, 0.6);
}

/* ========================================
   Dialog
   ======================================== */

.dialog-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.6);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 9999;
}

.dialog-box {
    background: #2a2a35;
    border-radius: 12px;
    padding: 24px;
    max-width: 320px;
    text-align: center;
    border: 1px solid rgba(255, 255, 255, 0.1);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
}

.dialog-icon {
    font-size: 40px;
    margin-bottom: 12px;
}

.dialog-box h4 {
    margin: 0 0 8px 0;
    font-size: 16px;
    color: #fff;
}

.dialog-box p {
    margin: 0 0 20px 0;
    font-size: 13px;
    color: rgba(255, 255, 255, 0.6);
    line-height: 1.5;
}

.dialog-actions {
    display: flex;
    gap: 10px;
    justify-content: center;
}

.btn-cancel,
.btn-confirm {
    padding: 8px 20px;
    border-radius: 6px;
    font-size: 13px;
    cursor: pointer;
    transition: all 0.2s;
}

.btn-cancel {
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.15);
    color: #fff;
}

.btn-cancel:hover {
    background: rgba(255, 255, 255, 0.15);
}

.btn-confirm {
    background: linear-gradient(135deg, #6366f1, #8b5cf6);
    border: none;
    color: #fff;
}

.btn-confirm:hover {
    box-shadow: 0 2px 12px rgba(99, 102, 241, 0.4);
}

/* ========================================
   Utilities
   ======================================== */

.spin {
    display: inline-block;
    animation: spin 1s linear infinite;
}

.spin.large {
    font-size: 32px;
}

@keyframes spin {
    from {
        transform: rotate(0deg);
    }

    to {
        transform: rotate(360deg);
    }
}
</style>
