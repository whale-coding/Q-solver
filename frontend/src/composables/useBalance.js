import { ref } from 'vue'

/**
 * 余额管理 composable
 * 注意：由于 OpenAI API 不提供标准的余额查询接口，这个 composable 已简化
 * 现在主要用于在设置面板中显示临时状态
 */
export function useBalance(settings, statusText, statusIcon, resetStatus) {
  const balance = ref(null)
  const tempBalance = ref(null)
  const isRefreshingBalance = ref(false)

  // 简化后不再主动查询余额
  async function fetchBalance(force = false) {
    // 不再实现，连通性测试移到模型选择页面
  }

  function refreshBalance() {
    // 不再实现
  }

  return {
    balance,
    tempBalance,
    isRefreshingBalance,
    fetchBalance,
    refreshBalance
  }
}
