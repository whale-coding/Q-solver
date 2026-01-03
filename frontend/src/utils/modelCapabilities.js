/**
 * 模型能力查询工具
 * 提供模型能力检测、厂商识别等功能
 */

import modelConfig from '../config/model-capabilities.json'

/**
 * 厂商信息配置
 * logo 使用 SVG data URI，确保图标清晰且无外部依赖
 */
export const PROVIDERS = {
    google: {
        name: 'Google',
        color: '#4285F4',
        // Google "G" logo
        logo: `<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"/>
      <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
      <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/>
      <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
    </svg>`
    },
    openai: {
        name: 'OpenAI',
        color: '#10A37F',
        // OpenAI logo
        logo: `<svg viewBox="0 0 24 24" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
      <path d="M22.2819 9.8211a5.9847 5.9847 0 0 0-.5157-4.9108 6.0462 6.0462 0 0 0-6.5098-2.9A6.0651 6.0651 0 0 0 4.9807 4.1818a5.9847 5.9847 0 0 0-3.9977 2.9 6.0462 6.0462 0 0 0 .7427 7.0966 5.98 5.98 0 0 0 .511 4.9107 6.051 6.051 0 0 0 6.5146 2.9001A5.9847 5.9847 0 0 0 13.2599 24a6.0557 6.0557 0 0 0 5.7718-4.2058 5.9894 5.9894 0 0 0 3.9977-2.9001 6.0557 6.0557 0 0 0-.7475-7.0729zm-9.022 12.6081a4.4755 4.4755 0 0 1-2.8764-1.0408l.1419-.0804 4.7783-2.7582a.7948.7948 0 0 0 .3927-.6813v-6.7369l2.02 1.1686a.071.071 0 0 1 .038.052v5.5826a4.504 4.504 0 0 1-4.4945 4.4944zm-9.6607-4.1254a4.4708 4.4708 0 0 1-.5346-3.0137l.142.0852 4.783 2.7582a.7712.7712 0 0 0 .7806 0l5.8428-3.3685v2.3324a.0804.0804 0 0 1-.0332.0615L9.74 19.9502a4.4992 4.4992 0 0 1-6.1408-1.6464zM2.3408 7.8956a4.485 4.485 0 0 1 2.3655-1.9728V11.6a.7664.7664 0 0 0 .3879.6765l5.8144 3.3543-2.0201 1.1685a.0757.0757 0 0 1-.071 0l-4.8303-2.7865A4.504 4.504 0 0 1 2.3408 7.8956zm16.5963 3.8558L13.1038 8.364 15.1192 7.2a.0757.0757 0 0 1 .071 0l4.8303 2.7913a4.4944 4.4944 0 0 1-.6765 8.1042v-5.6772a.79.79 0 0 0-.407-.667zm2.0107-3.0231l-.142-.0852-4.7735-2.7818a.7759.7759 0 0 0-.7854 0L9.409 9.2297V6.8974a.0662.0662 0 0 1 .0284-.0615l4.8303-2.7866a4.4992 4.4992 0 0 1 6.6802 4.66zM8.3065 12.863l-2.02-1.1638a.0804.0804 0 0 1-.038-.0567V6.0742a4.4992 4.4992 0 0 1 7.3757-3.4537l-.142.0805L8.704 5.459a.7948.7948 0 0 0-.3927.6813zm1.0976-2.3654l2.602-1.4998 2.6069 1.4998v2.9994l-2.5974 1.4997-2.6067-1.4997Z" fill="#10A37F"/>
    </svg>`
    },
    anthropic: {
        name: 'Anthropic',
        color: '#D97706',
        // Anthropic "A" style logo
        logo: `<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M17.304 3H13.896L20.4 21H23.808L17.304 3Z" fill="#D97706"/>
      <path d="M6.696 3L0.192 21H3.648L4.992 17.208H11.808L13.152 21H16.608L10.104 3H6.696ZM5.952 14.28L8.4 6.816L10.848 14.28H5.952Z" fill="#D97706"/>
    </svg>`
    },
    alibaba: {
        name: 'Alibaba',
        color: '#FF6A00',
        // Qwen/Alibaba Cloud style
        logo: `<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <circle cx="12" cy="12" r="10" fill="#FF6A00"/>
      <text x="12" y="16" text-anchor="middle" fill="white" font-size="10" font-weight="bold">Q</text>
    </svg>`
    },
    deepseek: {
        name: 'DeepSeek',
        color: '#0066FF',
        // DeepSeek style
        logo: `<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <circle cx="12" cy="12" r="10" fill="#0066FF"/>
      <path d="M8 12C8 9.79086 9.79086 8 12 8V8C14.2091 8 16 9.79086 16 12V12C16 14.2091 14.2091 16 12 16V16C9.79086 16 8 14.2091 8 12V12Z" stroke="white" stroke-width="1.5"/>
      <circle cx="12" cy="12" r="2" fill="white"/>
    </svg>`
    },
    zhipu: {
        name: '智谱AI',
        color: '#6366F1',
        logo: `<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect x="2" y="2" width="20" height="20" rx="4" fill="#6366F1"/>
      <text x="12" y="16" text-anchor="middle" fill="white" font-size="10" font-weight="bold">智</text>
    </svg>`
    },
    moonshot: {
        name: 'Moonshot',
        color: '#1a1a2e',
        // Moon icon
        logo: `<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <circle cx="12" cy="12" r="10" fill="#1a1a2e"/>
      <path d="M12 4C8.68629 4 6 7.58172 6 12C6 16.4183 8.68629 20 12 20C10.3431 20 9 16.4183 9 12C9 7.58172 10.3431 4 12 4Z" fill="#FFD700"/>
    </svg>`
    },
    '01ai': {
        name: '01.AI',
        color: '#EC4899',
        logo: `<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <circle cx="12" cy="12" r="10" fill="#EC4899"/>
      <text x="12" y="16" text-anchor="middle" fill="white" font-size="8" font-weight="bold">01</text>
    </svg>`
    },
    meta: {
        name: 'Meta',
        color: '#0668E1',
        // Meta infinity style
        logo: `<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M12 10.5C9.5 7 7 5.5 5 5.5C2.5 5.5 1 8 1 12C1 16 2.5 18.5 5 18.5C7 18.5 9.5 17 12 13.5C14.5 17 17 18.5 19 18.5C21.5 18.5 23 16 23 12C23 8 21.5 5.5 19 5.5C17 5.5 14.5 7 12 10.5Z" stroke="#0668E1" stroke-width="2" fill="none"/>
    </svg>`
    },
    mistral: {
        name: 'Mistral',
        color: '#F97316',
        logo: `<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect x="2" y="4" width="4" height="4" fill="#F97316"/>
      <rect x="10" y="4" width="4" height="4" fill="#F97316"/>
      <rect x="18" y="4" width="4" height="4" fill="#F97316"/>
      <rect x="2" y="10" width="4" height="4" fill="#F97316"/>
      <rect x="6" y="10" width="4" height="4" fill="#1a1a1a"/>
      <rect x="10" y="10" width="4" height="4" fill="#F97316"/>
      <rect x="14" y="10" width="4" height="4" fill="#1a1a1a"/>
      <rect x="18" y="10" width="4" height="4" fill="#F97316"/>
      <rect x="2" y="16" width="4" height="4" fill="#F97316"/>
      <rect x="10" y="16" width="4" height="4" fill="#F97316"/>
      <rect x="18" y="16" width="4" height="4" fill="#F97316"/>
    </svg>`
    },
    unknown: {
        name: '未知',
        color: '#6B7280',
        logo: `<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <circle cx="12" cy="12" r="10" stroke="#6B7280" stroke-width="2" fill="none"/>
      <text x="12" y="16" text-anchor="middle" fill="#6B7280" font-size="12" font-weight="bold">?</text>
    </svg>`
    }
}

/**
 * 能力图标配置
 */
export const CAPABILITY_ICONS = {
    text: { icon: '💬', label: '文本' },
    image: { icon: '🖼️', label: '图片' },
    pdf: { icon: '📄', label: 'PDF' },
    audio: { icon: '🎵', label: '音频' }
}

/**
 * 根据模型名称获取厂商 ID
 */
export function getProviderId(modelName) {
    if (!modelName) return 'unknown'

    const lower = modelName.toLowerCase()

    // 精确匹配
    const modelInfo = modelConfig.models[lower]
    if (modelInfo?.provider) {
        return modelInfo.provider
    }

    // 模式匹配
    for (const [pattern, provider] of Object.entries(modelConfig.providerPatterns)) {
        if (lower.includes(pattern.toLowerCase())) {
            return provider
        }
    }

    return 'unknown'
}

/**
 * 获取厂商信息
 */
export function getProvider(modelName) {
    const providerId = getProviderId(modelName)
    return PROVIDERS[providerId] || PROVIDERS.unknown
}

/**
 * 获取厂商 Logo SVG
 */
export function getProviderLogo(modelName) {
    return getProvider(modelName).logo
}

/**
 * 获取厂商名称
 */
export function getProviderName(modelName) {
    return getProvider(modelName).name
}

/**
 * 获取模型能力
 */
export function getModelCapabilities(modelName) {
    const defaultCaps = { text: true, image: false, pdf: false, audio: false }
    if (!modelName) return defaultCaps

    const lower = modelName.toLowerCase()

    // 精确匹配
    if (modelConfig.models[lower]) {
        const { provider, ...caps } = modelConfig.models[lower]
        return caps
    }

    // 模糊匹配
    for (const [key, value] of Object.entries(modelConfig.models)) {
        if (lower.includes(key) || key.includes(lower)) {
            const { provider, ...caps } = value
            return caps
        }
    }

    // 基于厂商的默认能力推断
    const providerId = getProviderId(modelName)
    if (providerId === 'google') {
        return { text: true, image: true, pdf: true, audio: false }
    }
    if (providerId === 'anthropic') {
        return { text: true, image: true, pdf: true, audio: false }
    }

    return defaultCaps
}

/**
 * 检查模型是否支持视觉（图片或PDF）
 */
export function supportsVision(modelName) {
    const caps = getModelCapabilities(modelName)
    return caps.image || caps.pdf
}

/**
 * 获取模型完整显示信息
 */
export function getModelInfo(modelName) {
    const provider = getProvider(modelName)
    const capabilities = getModelCapabilities(modelName)

    return {
        name: modelName,
        providerId: getProviderId(modelName),
        providerName: provider.name,
        providerColor: provider.color,
        providerLogo: provider.logo,
        capabilities,
        supportsVision: capabilities.image || capabilities.pdf
    }
}
