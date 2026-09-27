import dayjs from 'dayjs'

// 任务状态元数据：文案 + 主题色
// 数据来源为后端 model.Task 的状态常量
export const STATUS_META = {
  pending: { text: '等待中', color: '#f59e0b' },
  downloading: { text: '下载中', color: '#6366f1' },
  completed: { text: '已完成', color: '#10b981' },
  failed: { text: '下载失败', color: '#ef4444' },
  interrupted: { text: '已中断', color: '#f97316' }
}

export function statusMeta(status) {
  return STATUS_META[status] || { text: status || '未知', color: '#98a2b8' }
}

// 失败与中断都属于「异常结束」，用于行样式与错误提示
export function isErrorStatus(status) {
  return status === 'failed' || status === 'interrupted'
}

// 是否为「进行中/待处理」，即尚未结束
export function isActiveStatus(status) {
  return status === 'pending' || status === 'downloading'
}

// 下载耗时：已结束用 finished_at，进行中用当前时间
export function formatDuration(task) {
  if (!task || !task.created_at || task.status === 'pending') return '-'

  const start = dayjs(task.created_at)
  const end = task.finished_at ? dayjs(task.finished_at) : dayjs()
  const seconds = Math.max(0, end.diff(start, 'second'))

  if (seconds < 60) return `${seconds} 秒`

  const minutes = Math.floor(seconds / 60)
  const restSeconds = seconds % 60
  if (minutes < 60) return restSeconds ? `${minutes} 分 ${restSeconds} 秒` : `${minutes} 分`

  const hours = Math.floor(minutes / 60)
  const restMinutes = minutes % 60
  if (hours < 24) return restMinutes ? `${hours} 小时 ${restMinutes} 分` : `${hours} 小时`

  const days = Math.floor(hours / 24)
  const restHours = hours % 24
  return restHours ? `${days} 天 ${restHours} 小时` : `${days} 天`
}

// 进度副标题：速度与已下载/总大小（后端在下载中才会填充这些字段）
export function progressMeta(task) {
  if (!task || isErrorStatus(task.status)) return ''
  const parts = []
  if (task.speed) parts.push(task.speed)
  if (task.total_size) {
    parts.push(`${task.downloaded_size || '0'} / ${task.total_size}`)
  } else if (task.downloaded_size) {
    parts.push(task.downloaded_size)
  }
  return parts.join(' · ')
}
