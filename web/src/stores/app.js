import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'

const THEME_MODES = ['light', 'dark', 'system']

// 读取系统是否为暗色偏好
function systemPrefersDark() {
  return window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches
}

export const useAppStore = defineStore('app', () => {
  // 侧边栏折叠状态
  const collapsed = ref(localStorage.getItem('sidebarCollapsed') === 'true')

  // 主题模式：light / dark / system
  const theme = ref(
    THEME_MODES.includes(localStorage.getItem('theme')) ? localStorage.getItem('theme') : 'light'
  )

  // 系统偏好（system 模式下跟随它）
  const systemDark = ref(systemPrefersDark())

  // 实际生效的主题：system 模式下解析为当前系统偏好
  const resolvedTheme = computed(() =>
    theme.value === 'system' ? (systemDark.value ? 'dark' : 'light') : theme.value
  )

  // 监听变化并持久化
  watch(collapsed, (val) => {
    localStorage.setItem('sidebarCollapsed', val)
  })

  watch(theme, (val) => {
    localStorage.setItem('theme', val)
    applyTheme()
  })

  // 应用主题到 DOM
  function applyTheme() {
    document.documentElement.classList.toggle('dark', resolvedTheme.value === 'dark')
  }

  // 系统主题变化时实时响应（仅在 system 模式下有视觉效果）
  const mediaQuery = window.matchMedia ? window.matchMedia('(prefers-color-scheme: dark)') : null
  if (mediaQuery) {
    const onChange = (e) => {
      systemDark.value = e.matches
      if (theme.value === 'system') applyTheme()
    }
    // addEventListener 在旧版 Safari 上不可用，回退到 addListener
    if (mediaQuery.addEventListener) {
      mediaQuery.addEventListener('change', onChange)
    } else if (mediaQuery.addListener) {
      mediaQuery.addListener(onChange)
    }
  }

  // 切换侧边栏
  function toggleCollapsed() {
    collapsed.value = !collapsed.value
  }

  // 设置主题模式
  function setTheme(mode) {
    theme.value = THEME_MODES.includes(mode) ? mode : 'light'
  }

  // 在亮色/暗色之间快速切换（跟随系统时切到与当前相反的固定主题）
  function toggleTheme() {
    setTheme(resolvedTheme.value === 'dark' ? 'light' : 'dark')
  }

  // 初始化时应用主题
  applyTheme()

  return {
    collapsed,
    theme,
    resolvedTheme,
    setTheme,
    toggleCollapsed,
    toggleTheme
  }
})
