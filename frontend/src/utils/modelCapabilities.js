import modelData from '../config/model-capabilities.json'

export const PROVIDERS = {
    openai: { name: 'OpenAI' },
    google: { name: 'Google', aliases: ['deepmind'] },
    anthropic: { name: 'Anthropic' },
    meta: { name: 'Meta', aliases: ['facebook'] },
    nvidia: { name: 'NVIDIA' },
    microsoft: { name: 'Microsoft', aliases: ['azure'] },
    deepseek: { name: 'DeepSeek' },
    zhipu: { name: 'Zhipu AI', aliases: ['z-ai', 'bigmodel', 'chatglm', 'thmffu'] },
    moonshot: { name: 'Moonshot AI', aliases: ['moonshotai', 'kimi'] },
    alibaba: { name: 'Alibaba Cloud', aliases: ['qwen'] },
    bytedance: { name: 'ByteDance', aliases: ['volcengine', 'doubao', 'byte-dance'] },
    '01-ai': { name: '01.AI', aliases: ['01ai', 'yi', 'zeroone'] },
    mistral: { name: 'Mistral AI', aliases: ['mistralai', 'mistral-ai'] },
    cohere: { name: 'Cohere' },
    perplexity: { name: 'Perplexity' },
    xai: { name: 'xAI', aliases: ['x-ai', 'grok'] },
    liquid: { name: 'Liquid AI' },
    upstage: { name: 'Upstage' },
    openrouter: { name: 'OpenRouter' },
    minimax: { name: 'MiniMax' },
    baichuan: { name: 'Baichuan' },
    huggingface: { name: 'HuggingFace', aliases: ['nousresearch', 'cognitivecomputations'] },
    together: { name: 'Together AI' },
    fireworks: { name: 'Fireworks AI' },
    allenai: { name: 'AllenAI', aliases: ['ai2'] },
    xiaomi: { name: 'Xiaomi', aliases: ['xiaomimimo'] },
    ibm: { name: 'IBM', aliases: ['ibm-granite'] },
    deepcogito: { name: 'Deep Cogito', aliases: ['deep-cogito'] },
    amazon: { name: 'Amazon', aliases: ['aws', 'bedrock'] },
    custom: { name: '自定义' },
}

// 默认 SVG 图标 (写死在代码里作为兜底)
const DEFAULT_LOGO_SVG = `<svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
    <rect width="24" height="24" rx="6" fill="#4B5563"/>
    <rect x="7" y="7" width="10" height="10" rx="2" stroke="white" stroke-width="2" fill="none"/>
    <circle cx="12" cy="12" r="2" fill="white"/>
</svg>`

export const PROVIDER_BASE_URLS = {
    google: 'https://generativelanguage.googleapis.com',
    openai: 'https://api.openai.com/v1',
    anthropic: 'https://api.anthropic.com',
    alibaba: 'https://dashscope.aliyuncs.com/compatible-mode/v1',
    moonshot: 'https://api.moonshot.cn/v1',
    openrouter: 'https://openrouter.ai/api/v1',
    custom: ''
}

// 默认能力
const defaultCapabilities = {
    text: true,
    image: false,
    file: false,
    audio: false,
    video: false,
    contextLength: 0
}

/**
 * 解析模型能力从 inputs 数组
 */
function parseCapabilities(modelInfo) {
    const inputs = modelInfo.inputs || []
    return {
        text: inputs.includes('text'),
        image: inputs.includes('image') || modelInfo.supports_image === true,
        file: inputs.includes('file'),
        audio: inputs.includes('audio'),
        video: inputs.includes('video'),
        contextLength: modelInfo.context_length || 0
    }
}

/**
 * 查找模型信息，支持精确匹配和模糊匹配
 */
function findModelInfo(model) {
    if (!model) return null

    // 1. 清理 key：去掉厂商前缀，转小写
    const cleanKey = model.toLowerCase().split('/').pop()

    // 2. 精确匹配
    if (modelData[cleanKey]) {
        return modelData[cleanKey]
    }

    // 3. 模糊匹配：遍历所有 key 查找包含关系
    for (const [key, value] of Object.entries(modelData)) {
        if (cleanKey.includes(key) || key.includes(cleanKey)) {
            return value
        }
    }

    return null
}

/**
 * 标准化 provider code - 处理别名映射
 */
function normalizeProviderCode(code) {
    if (!code) return 'unknown'
    const lowerCode = code.toLowerCase()

    // 精确匹配
    if (PROVIDERS[lowerCode]) {
        return lowerCode
    }

    // 检查别名
    for (const [key, value] of Object.entries(PROVIDERS)) {
        if (value.aliases?.includes(lowerCode)) {
            return key
        }
    }

    return lowerCode
}

/**
 * Returns the capabilities object for a given model
 */
export function getModelCapabilities(model) {
    if (!model) return defaultCapabilities

    const modelInfo = findModelInfo(model)
    if (modelInfo) {
        return parseCapabilities(modelInfo)
    }
    return defaultCapabilities
}

/**
 * 获取模型的完整信息
 */
export function getModelInfo(model) {
    return findModelInfo(model)
}

/**
 * Returns the friendly provider name for a given model
 */
export function getProviderName(model) {
    if (!model) return 'Unknown'

    const modelInfo = findModelInfo(model)
    if (modelInfo && modelInfo.provider) {
        const normalized = normalizeProviderCode(modelInfo.provider)
        if (PROVIDERS[normalized]) {
            return PROVIDERS[normalized].name
        }
        // 没找到就 capitalize
        return modelInfo.provider.charAt(0).toUpperCase() + modelInfo.provider.slice(1)
    }

    return 'Unknown'
}

/**
 * Returns the provider code for a given model
 */
export function getProviderCode(model) {
    if (!model) return 'unknown'

    const modelInfo = findModelInfo(model)
    if (modelInfo && modelInfo.provider) {
        return normalizeProviderCode(modelInfo.provider)
    }

    return 'unknown'
}

/**
 * Provider code 到实际图标文件名的映射
 * lobehub/icons 命名不统一，需要手动映射
 */
const ICON_FILENAME_MAP = {
    // 有 -color.svg 后缀的
    'google': 'gemini-color',
    'deepseek': 'deepseek-color',
    'bytedance': 'bytedance-color',
    'cohere': 'cohere-color',
    'perplexity': 'perplexity-color',
    'nvidia': 'nvidia-color',
    'huggingface': 'huggingface-color',
    'minimax': 'minimax-color',
    'baichuan': 'baichuan-color',
    'meta': 'meta-color',
    'deepcogito': 'deepcogito-color',
    'mistral': 'mistral-color',
    'fireworks': 'fireworks-color',
    'upstage': 'upstage-color',

    // 没有 -color.svg 后缀的，直接用原名
    'openai': 'openai',
    'anthropic': 'anthropic',
    'xai': 'grok',  // xAI 用 grok 图标
    'moonshot': 'moonshot',
    'alibaba': 'qwen-color',
    'qwen': 'qwen-color',
    'zhipu': 'zhipu-color',
    'openrouter': 'openrouter',
    'allenai': 'ai2-color',
    'ai2': 'ai2-color',
    '01-ai': 'yi-color',
    'yi': 'yi-color',
    'zeroone': 'zeroone',
    'amazon': 'aws-color',
    'aws': 'aws-color',
    'ibm': 'ibm',
    'ibm-granite': 'ibm',
    'liquid': 'liquid',
    'together': 'together-color',
    'xiaomi': 'xiaomimimo',
    'inflection': 'inflection',
    'ai21': 'ai21',
    'arcee': 'arcee-color',

    // 常见别名映射
    'meta-llama': 'meta-color',
    'bytedance-seed': 'bytedance-color',
    'stepfun': 'stepfun-color',
    'stepfun-ai': 'stepfun-color',
    'microsoft': 'microsoft-color',
    'azure': 'azure-color',
    'nousresearch': 'nousresearch',
    'sambanova': 'sambanova-color',
    'cerebras': 'cerebras-color',
    'groq': 'groq',
    'arcee-ai': 'arcee-color',
    'x-ai': 'grok',
    'aion-labs': 'aionlabs-color',

    // 小众/社区模型厂商 - 没有专用图标，使用 null 触发默认图标
    'eleutherai': null,
    'alpindale': null,
    'raifle': null,
    'anthracite-org': null,
    'nex-agi': null,
    'prime-intellect': null,
    'tngtech': null,
    'thedrummer': null,
    'meituan': null,
    'opengvlab': null,
    'switchpoint': null,
    'alfredpros': null,
    'sao10k': null,
    'neversleep': null,
    'undi95': null,
    'mancer': null,
    'gryphe': null,
    'kwaipilot': null,
    'essentialai': null,
    'relace': null,
}

/**
 * Returns the logo path/SVG for a given model's provider
 * 优先使用 /icons/ (lobehub/icons), 兜底返回内联 SVG
 */
export function getProviderLogo(model) {
    const providerCode = getProviderCode(model)

    if (providerCode === 'unknown' || providerCode === 'custom') {
        return DEFAULT_LOGO_SVG
    }

    // O(1) 查映射表
    const iconName = ICON_FILENAME_MAP[providerCode]
    if (iconName !== undefined) {
        // null 表示没有专用图标，使用默认图标
        return iconName === null ? DEFAULT_LOGO_SVG : `/icons/${iconName}.svg`
    }

    // 没在映射表中，直接用 provider code 作为文件名
    return `/icons/${providerCode}.svg`
}

/**
 * 获取默认 logo SVG (供外部使用)
 */
export function getDefaultLogoSvg() {
    return DEFAULT_LOGO_SVG
}

/**
 * 格式化上下文长度
 */
export function formatContextLength(length) {
    if (!length || length <= 0) return ''

    if (length >= 1000000) {
        const m = length / 1000000
        return m % 1 === 0 ? `${m}M` : `${m.toFixed(1)}M`
    }

    if (length >= 1000) {
        const k = length / 1000
        return k % 1 === 0 ? `${k}K` : `${k.toFixed(0)}K`
    }

    return String(length)
}

/**
 * Check if the model supports vision capabilities
 */
export function supportsVision(model) {
    const capabilities = getModelCapabilities(model)
    return !!capabilities.image
}

/**
 * Check if the model supports file/PDF capabilities
 */
export function supportsFile(model) {
    const capabilities = getModelCapabilities(model)
    return !!capabilities.file
}

// Backward compatibility aliases
export const supportsPDF = supportsFile

// Legacy exports (保持向后兼容)
export const PROVIDER_LOGOS = new Proxy({}, {
    get(_, key) {
        if (key === 'default') return DEFAULT_LOGO_SVG
        return `/icons/${key}.svg`
    }
})

export const PROVIDER_NAMES = new Proxy({}, {
    get(_, key) {
        const normalized = normalizeProviderCode(key)
        return PROVIDERS[normalized]?.name || key
    }
})
