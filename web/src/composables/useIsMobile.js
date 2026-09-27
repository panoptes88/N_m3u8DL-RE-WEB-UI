import { ref, onMounted, onUnmounted } from 'vue'

const MOBILE_MAX_WIDTH = 768

/**
 * 响应式判断当前是否为移动端宽度（<=768px）。
 * 组件的挂载/卸载会自动注册与移除监听。
 */
export function useIsMobile() {
  const isMobile = ref(false)

  const update = () => {
    isMobile.value = window.innerWidth <= MOBILE_MAX_WIDTH
  }

  // 立即求值，避免首帧先按桌面渲染再跳变
  update()

  onMounted(() => window.addEventListener('resize', update))
  onUnmounted(() => window.removeEventListener('resize', update))

  return { isMobile }
}
