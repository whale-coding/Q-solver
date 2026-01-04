<template>
    <div class="resume-import-container">
        <!-- Header / File Selection -->
        <div class="resume-header-section">
            <div class="header-title">
                <span class="icon">📄</span>
                <span>简历导入</span>
            </div>

            <div class="file-controls">
                <div class="toggle-wrapper" v-if="resumePath"
                    @click="$emit('update:useMarkdownResume', !useMarkdownResume)"
                    title="勾选后，AI 将使用解析后的 Markdown 文本作为简历。不勾选则发送 PDF 版本。">
                    <span class="toggle-label" :class="{ active: useMarkdownResume }">使用MarkDown简历</span>
                    <div class="toggle-switch" :class="{ active: useMarkdownResume }">
                        <div class="toggle-knob"></div>
                    </div>
                </div>
                <div class="separator-v" v-if="resumePath"></div>

                <div v-if="!resumePath" class="upload-btn" @click="$emit('select-resume')">
                    <span class="icon">📂</span> 选择 PDF 简历
                </div>
                <div v-else class="file-selected">
                    <div class="btn-group">
                        <button class="btn-secondary small" @click="toggleFlip">
                            {{ isFlipped ? '切换 PDF' : '切换 Markdown' }}
                        </button>
                        <button class="btn-secondary small" @click="enableManualInput">手动输入</button>
                        <button class="btn-secondary small" @click="$emit('select-resume')">更换</button>
                        <button class="btn-danger small" @click="$emit('clear-resume')">清除</button>
                    </div>
                </div>
            </div>
        </div>

        <!-- Flip Container -->
        <div class="flip-container" v-if="resumePath">
            <div class="flipper" :class="{ flipped: isFlipped }">
                <!-- Front: PDF Preview -->
                <div class="front">
                    <div class="preview-pane pdf-pane">
                        <div class="pane-label pdf-pane-label">
                            <span>PDF 预览</span>
                            <div class="pdf-pane-actions unified-actions">
                                <button class="btn-primary parse-btn" @click="handleParseClick" :disabled="isParsing">
                                    <span class="icon" v-if="!isParsing">✨</span>
                                    <span class="icon spin" v-else>⏳</span>
                                    {{ isParsing ? '解析中...' : 'AI 解析为 Markdown' }}
                                </button>
                                <span v-if="!modelSupportsFile" class="vision-warning" title="当前模型可能不支持 PDF，点击解析时会提示确认">
                                    ⚠️
                                </span>
                                <template v-if="!pdfControlsCollapsed">
                                    <button class="pdf-btn" @click="prevPage" :disabled="pageNum <= 1" title="上一页">
                                        &lt;
                                    </button>
                                    <span class="page-info">{{ pageNum }} / {{ pageCount }}</span>
                                    <button class="pdf-btn" @click="nextPage" :disabled="pageNum >= pageCount"
                                        title="下一页">
                                        &gt;
                                    </button>
                                </template>
                            </div>
                        </div>
                        <div class="pane-body">
                            <div v-show="pageCount > 0" class="pdf-container">
                                <div class="canvas-wrapper">
                                    <canvas ref="canvasRef"></canvas>
                                </div>
                            </div>
                            <div v-if="pageCount === 0" class="placeholder-content">
                                <span class="icon">📑</span>
                                <p>PDF 预览区域</p>
                                <p class="sub-text">{{ resumePath }}</p>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Back: Markdown Preview -->
                <div class="back">
                    <div class="preview-pane markdown-pane">
                        <div class="pane-label">
                            <span>Markdown 预览</span>
                            <div class="right-controls">
                                <span class="action-text" @click="toggleEdit" v-if="renderedContent || localContent">{{
                                    isEditing ? '预览'
                                        : '编辑' }}</span>
                                <span class="separator" v-if="renderedContent || localContent">|</span>
                                <span class="hint-text" @click="toggleFlip">&lt; 点击切换到 PDF</span>
                            </div>
                        </div>
                        <div class="pane-body">
                            <div v-if="isEditing" class="editor-container">
                                <textarea v-model="localContent" @input="updateContent" class="markdown-editor"
                                    placeholder="在此编辑 Markdown..."></textarea>
                            </div>
                            <div v-else-if="renderedContent" class="markdown-content" v-html="renderedContent"></div>
                            <div v-else-if="isParsing" class="placeholder-content">
                                <span class="icon spin">⏳</span>
                                <p>正在解析中...</p>
                            </div>
                            <div v-else class="placeholder-content">
                                <span class="icon">📝</span>
                                <p>等待解析...</p>
                            </div>
                        </div>
                    </div>
                </div>

            </div>
        </div>

        <!-- Empty State Description -->
        <div v-else class="empty-state-desc">
            <p>导入您的 PDF 简历，AI 将在解题时参考您的背景信息，提供更个性化的回答。</p>
        </div>

        <!-- Bottom Action Bar 已移除 -->

        <!-- PDF 解析确认弹窗 -->
        <div v-if="showConfirmDialog" class="confirm-overlay">
            <div class="confirm-dialog">
                <div class="confirm-icon">⚠️</div>
                <div class="confirm-title">模型可能不支持</div>
                <div class="confirm-message">当前模型可能不支持 PDF 解析，是否仍要继续？</div>
                <div class="confirm-actions">
                    <button class="btn-secondary" @click="showConfirmDialog = false">取消</button>
                    <button class="btn-primary" @click="confirmParse">继续</button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { computed, ref, watch, onMounted, nextTick } from 'vue';
import { marked } from 'marked';
import { GetResumePDF } from '../../wailsjs/go/main/App';
import * as pdfjsLib from 'pdfjs-dist';
import { supportsVision, supportsPDF } from '../utils/modelCapabilities';

pdfjsLib.GlobalWorkerOptions.workerSrc = new URL(
    'pdfjs-dist/build/pdf.worker.mjs',
    import.meta.url
).href;

const props = defineProps({
    resumePath: {
        type: String,
        default: ''
    },
    rawContent: {
        type: String,
        default: ''
    },
    isParsing: {
        type: Boolean,
        default: false
    },
    useMarkdownResume: {
        type: Boolean,
        default: false
    },
    currentModel: {
        type: String,
        default: ''
    }
});

// 检查当前模型是否支持视觉或 PDF 功能
const modelSupportsFile = computed(() => supportsVision(props.currentModel) || supportsPDF(props.currentModel));

const emit = defineEmits(['select-resume', 'clear-resume', 'parse-resume', 'update:rawContent', 'update:useMarkdownResume']);

const isFlipped = ref(false);
const isEditing = ref(false);
const localContent = ref(props.rawContent);
const showConfirmDialog = ref(false);

// PDF 相关状态
const pageNum = ref(1);
const pageCount = ref(0);
const scale = ref(0.75);
const canvasRef = ref(null);
let pdfDoc = null;
let renderTask = null;

watch(() => props.rawContent, (newVal) => {
    if (newVal !== localContent.value) {
        localContent.value = newVal;
    }
});

watch(() => props.resumePath, async (newVal) => {
    if (newVal) {
        await loadPdfPreview();
    } else {
        pdfDoc = null;
        pageCount.value = 0;
        pageNum.value = 1;
        // 清除 canvas
        const canvas = canvasRef.value;
        if (canvas) {
            const ctx = canvas.getContext('2d');
            ctx.clearRect(0, 0, canvas.width, canvas.height);
        }
    }
});

onMounted(async () => {
    if (props.resumePath) {
        await loadPdfPreview();
    }
});

async function loadPdfPreview() {
    try {
        const base64 = await GetResumePDF();
        if (base64) {
            const binaryString = window.atob(base64);
            const len = binaryString.length;
            const bytes = new Uint8Array(len);
            for (let i = 0; i < len; i++) {
                bytes[i] = binaryString.charCodeAt(i);
            }

            const loadingTask = pdfjsLib.getDocument({ data: bytes });
            pdfDoc = await loadingTask.promise;
            pageCount.value = pdfDoc.numPages;
            pageNum.value = 1;
            // 确保 DOM 更新后渲染
            nextTick(() => renderPage(pageNum.value));
        }
    } catch (e) {
        console.error("Failed to load PDF preview:", e);
    }
}

async function renderPage(num) {
    if (!pdfDoc) return;

    try {
        const page = await pdfDoc.getPage(num);
        const canvas = canvasRef.value;
        if (!canvas) return;

        // 如果有正在进行的渲染任务，取消它
        if (renderTask) {
            renderTask.cancel();
        }

        const ctx = canvas.getContext('2d');
        const viewport = page.getViewport({ scale: scale.value });

        canvas.height = viewport.height;
        canvas.width = viewport.width;

        const renderContext = {
            canvasContext: ctx,
            viewport: viewport
        };

        renderTask = page.render(renderContext);
        await renderTask.promise;
    } catch (e) {
        if (e.name !== 'RenderingCancelledException') {
            console.error("Render error:", e);
        }
    } finally {
        renderTask = null;
    }
}

function prevPage() {
    if (pageNum.value <= 1) return;
    pageNum.value--;
    renderPage(pageNum.value);
}

function nextPage() {
    if (pageNum.value >= pageCount.value) return;
    pageNum.value++;
    renderPage(pageNum.value);
}

function zoomIn() {
    scale.value += 0.25;
    renderPage(pageNum.value);
}

function zoomOut() {
    if (scale.value > 0.5) {
        scale.value -= 0.25;
        renderPage(pageNum.value);
    }
}

function updateContent() {
    emit('update:rawContent', localContent.value);
}

// 使用 marked.parse 解析 markdown 内容
const renderedContent = computed(() => {
    if (!localContent.value) return '';
    return marked.parse(localContent.value);
});

function toggleFlip() {
    isFlipped.value = !isFlipped.value;
}

function toggleEdit() {
    isEditing.value = !isEditing.value;
}

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

function enableManualInput() {
    isFlipped.value = true;
    isEditing.value = true;
}
</script>


<style scoped>
.spin {
    animation: spin 1s linear infinite;
    display: inline-block;
}

@keyframes spin {
    from {
        transform: rotate(0deg);
    }

    to {
        transform: rotate(360deg);
    }
}

.markdown-content {
    flex: 1;
    width: 100%;
    height: 100%;
    overflow-y: auto;
    overflow-x: hidden;
    padding: 15px;
    text-align: left;
    font-size: 13px;
    line-height: 1.6;
    color: #e0e0e0;
    min-height: 0;
    /* 解决 3D 变换导致滚动失效的问题 */
    transform: translateZ(0);
    -webkit-overflow-scrolling: touch;
}

.markdown-content :deep(h1),
.markdown-content :deep(h2),
.markdown-content :deep(h3),
.markdown-content :deep(h4),
.markdown-content :deep(h5),
.markdown-content :deep(h6) {
    margin: 0.8em 0 0.4em 0;
    font-weight: 600;
}

.markdown-content :deep(h1) {
    font-size: 1.5em;
}

.markdown-content :deep(h2) {
    font-size: 1.3em;
}

.markdown-content :deep(h3) {
    font-size: 1.1em;
}

.markdown-content :deep(p) {
    margin: 0.5em 0;
}

.markdown-content :deep(ul),
.markdown-content :deep(ol) {
    margin: 0.5em 0;
    padding-left: 1.5em;
}

.markdown-content :deep(li) {
    margin: 0.25em 0;
}

.markdown-content :deep(code) {
    background: rgba(255, 255, 255, 0.1);
    padding: 0.15em 0.4em;
    border-radius: 3px;
    font-family: 'Fira Code', monospace;
    font-size: 0.9em;
}

.markdown-content :deep(pre) {
    background: rgba(0, 0, 0, 0.3);
    padding: 10px;
    border-radius: 5px;
    overflow-x: auto;
}

.markdown-content :deep(pre code) {
    background: transparent;
    padding: 0;
}

.markdown-content :deep(blockquote) {
    border-left: 3px solid #1890ff;
    margin: 0.5em 0;
    padding-left: 1em;
    color: #aaa;
}

.markdown-content :deep(hr) {
    border: none;
    border-top: 1px solid rgba(255, 255, 255, 0.1);
    margin: 1em 0;
}

.markdown-content :deep(a) {
    color: #1890ff;
}

.markdown-content :deep(strong) {
    font-weight: 600;
}

.markdown-content :deep(table) {
    width: 100%;
    border-collapse: collapse;
    margin: 0.5em 0;
}

.markdown-content :deep(th),
.markdown-content :deep(td) {
    border: 1px solid rgba(255, 255, 255, 0.1);
    padding: 0.5em;
    text-align: left;
}

.markdown-content :deep(th) {
    background: rgba(255, 255, 255, 0.05);
}

.markdown-editor {
    flex: 1;
    width: 100%;
    height: 100%;
    background: transparent;
    border: none;
    color: #e0e0e0;
    font-family: 'Fira Code', monospace;
    font-size: 13px;
    resize: none;
    outline: none;
    line-height: 1.6;
    padding: 15px;
    overflow-y: auto;
    transform: translateZ(0);
    -webkit-overflow-scrolling: touch;
}

.resume-import-container {
    display: flex;
    flex-direction: column;
    height: 100%;
    gap: 15px;
    color: #fff;
}

.resume-header-section {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-bottom: 10px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.header-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 16px;
    font-weight: 600;
}

.file-controls {
    display: flex;
    align-items: center;
}

.upload-btn {
    cursor: pointer;
    background: rgba(255, 255, 255, 0.1);
    padding: 6px 12px;
    border-radius: 6px;
    font-size: 13px;
    transition: background 0.2s;
    display: flex;
    align-items: center;
    gap: 6px;
}

.upload-btn:hover {
    background: rgba(255, 255, 255, 0.2);
}

.file-selected {
    display: flex;
    align-items: center;
    gap: 10px;
}

.btn-group {
    display: flex;
    gap: 5px;
}

.btn-secondary,
.btn-danger {
    padding: 4px 8px;
    border-radius: 4px;
    border: none;
    cursor: pointer;
    font-size: 12px;
    color: #fff;
}

.btn-danger {
    background: #ff4d4f;
    color: #fff;
    border: none;
}

.btn-danger:hover {
    background: #ff7875;
}

.btn-secondary {
    background: rgba(255, 255, 255, 0.15);
}

.btn-secondary:hover {
    background: rgba(255, 255, 255, 0.25);
}

/* Flip Animation Styles */
.flip-container {
    flex: 1;
    perspective: 1000px;
    min-height: 500px;
    /* 增加最小高度，让可视区域更大 */
    position: relative;
}

btn-danger:hover {
    background: rgba(255, 77, 79, 0.3);
}

/* Flip Animation Styles */
.flip-container {
    flex: 1;
    perspective: 1000px;
    min-height: 400px;
    position: relative;
}

.flipper {
    transition: 0.6s;
    transform-style: preserve-3d;
    position: relative;
    width: 100%;
    height: 100%;
}

.flipper.flipped {
    transform: rotateY(180deg);
}

.front,
.back {
    backface-visibility: hidden;
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
}

.front {
    z-index: 2;
    transform: rotateY(0deg);
}

.back {
    transform: rotateY(180deg);
    overflow: auto;
}

/* Pane Styles */
.preview-pane {
    height: 100%;
    background: rgba(0, 0, 0, 0.2);
    border-radius: 8px;
    border: 1px solid rgba(255, 255, 255, 0.05);
    display: flex;
    flex-direction: column;
    overflow: auto;
}

.pane-label {
    padding: 8px 12px;
    background: rgba(255, 255, 255, 0.05);
    font-size: 12px;
    color: #aaa;
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.pdf-pane-label {
    position: relative;
}

.pdf-pane-actions {
    display: flex;
    align-items: center;
    gap: 8px;
}

.parse-btn {
    font-size: 12px;
    padding: 4px 10px;
    border-radius: 4px;
    min-width: unset;
    background: rgba(24, 144, 255, 0.12);
    color: #1890ff;
    border: 1px solid #1890ff;
    box-shadow: none;
    transition: background 0.2s, color 0.2s, border 0.2s;
    font-weight: 500;
}

.parse-btn:disabled {
    background: rgba(24, 144, 255, 0.08);
    color: #aaa;
    border-color: #aaa;
    cursor: not-allowed;
}

.parse-btn:hover:not(:disabled) {
    background: #1890ff;
    color: #fff;
}

.vision-warning {
    font-size: 14px;
    cursor: help;
    animation: pulse 2s infinite;
}

@keyframes pulse {

    0%,
    100% {
        opacity: 1;
    }

    50% {
        opacity: 0.5;
    }
}


.hint-text {
    cursor: pointer;
    color: #1890ff;
    font-size: 11px;
}

.hint-text:hover {
    text-decoration: underline;
}

.pane-body {
    flex: 1;
    display: flex;
    flex-direction: column;
    position: relative;
    overflow: auto;
    background: #1e1e1e;
    min-height: 0;
}

.separator-v-small {
    width: 1px;
    height: 16px;
    background: rgba(255, 255, 255, 0.1);
    margin: 0 5px;
}

.icon-text {
    font-size: 16px;
    font-weight: bold;
    line-height: 1;
}

.pdf-container {
    width: 100%;
    height: 100%;
    overflow: hidden;
    background: #2d2d2d;
    display: flex;
    flex-direction: column;
}

.pdf-controls {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 10px;
    padding: 8px;
    background: rgba(0, 0, 0, 0.3);
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);
    z-index: 10;
    flex-shrink: 0;
    max-height: 44px;
    overflow: hidden;
    transition: max-height 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.pdf-controls.collapsed {
    max-height: 0;
    padding: 0 8px;
    border-bottom: none;
}

.collapse-btn {
    width: 28px;
    height: 28px;
    background: rgba(255, 255, 255, 0.1);
    border: none;
    color: #fff;
    border-radius: 4px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 16px;
    font-weight: bold;
    margin-right: 4px;
}

.canvas-wrapper {
    flex: 1;
    overflow: auto;
    display: flex;
    justify-content: center;
    align-items: center;
    width: auto;
    padding: 20px;
    box-sizing: border-box;
    background: #2d2d2d;
    position: relative;
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.pdf-btn {
    background: rgba(255, 255, 255, 0.1);
    border: none;
    color: #fff;
    width: 28px;
    height: 28px;
    border-radius: 4px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background 0.2s;
}

.pdf-btn:hover:not(:disabled) {
    background: rgba(255, 255, 255, 0.2);
}

.pdf-btn:disabled {
    opacity: 0.3;
    cursor: not-allowed;
}

.page-info {
    font-size: 12px;
    color: #ccc;
    font-variant-numeric: tabular-nums;
}

canvas {
    /* 阴影效果 */
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.5);

    /* 强制重置尺寸限制，由 JS 控制大小 */
    max-width: none !important;
    width: auto !important;
    height: auto !important;

    /* 布局关键：变成块级 + 自动边距实现居中 */
    display: block;
    margin: auto;
}

.placeholder-content {
    text-align: center;
    color: rgba(255, 255, 255, 0.3);
    margin: auto;
    padding: 20px;
}

.placeholder-content .icon {
    font-size: 32px;
    display: block;
    margin-bottom: 10px;
}

.sub-text {
    font-size: 11px;
    margin-top: 5px;
    opacity: 0.7;
    word-break: break-all;
}

.empty-state-desc {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    text-align: center;
    color: #888;
    padding: 20px;
    font-size: 14px;
    line-height: 1.5;
}

.action-bar {
    padding-top: 10px;
}

.btn-primary {
    background: #0986fc;
    color: white;
    border: none;
    /* padding: 10px 20px; */
    border-radius: 6px;
    cursor: pointer;
    font-size: 13px;
    font-weight: 700;
    transition: background 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
}

.btn-primary:hover {
    background: #40a9ff;
}

.full-width {
    width: 100%;
}

.right-controls {
    display: flex;
    align-items: center;
    gap: 8px;
}

.separator-v {
    width: 1px;
    height: 16px;
    background: rgba(255, 255, 255, 0.15);
    margin: 0 12px;
}

.toggle-wrapper {
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    user-select: none;
}

.toggle-label {
    font-size: 12px;
    color: #888;
    transition: color 0.3s ease;
    font-weight: 500;
}

.toggle-label.active {
    color: #fff;
}

.toggle-switch {
    width: 32px;
    height: 18px;
    background: rgba(255, 255, 255, 0.15);
    border-radius: 10px;
    position: relative;
    transition: all 0.3s cubic-bezier(0.4, 0.0, 0.2, 1);
    border: 1px solid rgba(255, 255, 255, 0.05);
}

.toggle-switch.active {
    background: #1890ff;
    border-color: #1890ff;
}

.toggle-knob {
    width: 14px;
    height: 14px;
    background: #fff;
    border-radius: 50%;
    position: absolute;
    top: 1px;
    left: 1px;
    transition: transform 0.3s cubic-bezier(0.4, 0.0, 0.2, 1);
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
}

.toggle-switch.active .toggle-knob {
    transform: translateX(14px);
}

.checkbox-wrapper {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 12px;
    color: #ccc;
    cursor: pointer;
}

.checkbox-wrapper input {
    cursor: pointer;
}

.checkbox-wrapper label {
    cursor: pointer;
}

.action-text {
    cursor: pointer;
    color: #1890ff;
    font-size: 12px;
}

.action-text:hover {
    text-decoration: underline;
}

.separator {
    color: rgba(255, 255, 255, 0.2);
}

.editor-container {
    flex: 1;
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    overflow: auto;
    min-height: 0;
}

/* Confirmation Dialog Styles */
.confirm-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.6);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
}

.confirm-dialog {
    background: #2a2a2a;
    border-radius: 12px;
    padding: 24px;
    max-width: 320px;
    text-align: center;
    border: 1px solid rgba(255, 255, 255, 0.1);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
}

.confirm-icon {
    font-size: 48px;
    margin-bottom: 12px;
}

.confirm-title {
    font-size: 16px;
    font-weight: 600;
    margin-bottom: 8px;
    color: #fff;
}

.confirm-message {
    font-size: 13px;
    color: #aaa;
    margin-bottom: 20px;
    line-height: 1.5;
}

.confirm-actions {
    display: flex;
    gap: 12px;
    justify-content: center;
}

.confirm-actions button {
    min-width: 80px;
    padding: 8px 16px;
    border-radius: 6px;
    font-size: 13px;
    cursor: pointer;
    transition: all 0.2s;
}

.confirm-actions .btn-primary {
    background: #1890ff;
    color: #fff;
    border: none;
}

.confirm-actions .btn-primary:hover {
    background: #40a9ff;
}

.confirm-actions .btn-secondary {
    background: rgba(255, 255, 255, 0.1);
    color: #fff;
    border: 1px solid rgba(255, 255, 255, 0.2);
}

.confirm-actions .btn-secondary:hover {
    background: rgba(255, 255, 255, 0.15);
}
</style>
