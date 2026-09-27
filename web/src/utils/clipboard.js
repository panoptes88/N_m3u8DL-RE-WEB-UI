import { message } from 'ant-design-vue'

// 复制到剪贴板
// 优先用 Clipboard API（需安全上下文：HTTPS 或 localhost），
// 不可用时回退到 execCommand（兼容 http://IP 访问）
export async function copyToClipboard(text) {
  if (!text) return false

  try {
    if (navigator.clipboard) {
      await navigator.clipboard.writeText(text)
      message.success('已复制到粘贴板')
      return true
    }
  } catch {
    // 忽略，回退到 execCommand
  }

  try {
    const textarea = document.createElement('textarea')
    textarea.value = text
    textarea.style.position = 'fixed'
    textarea.style.opacity = '0'
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
    message.success('已复制到粘贴板')
    return true
  } catch {
    message.error('复制失败')
    return false
  }
}
