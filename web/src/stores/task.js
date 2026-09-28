import { defineStore } from 'pinia'
import { ref } from 'vue'
import { get, post, del } from '../api'

export const useTaskStore = defineStore('task', () => {
  const tasks = ref([])
  const loading = ref(false)
  // 状态筛选收归 store 管理：轮询必须复用同一条件，
  // 否则筛选结果会在下一次轮询时被全量数据覆盖
  const statusFilter = ref('')
  let pollingTimer = null
  let polling = false

  // 自适应轮询间隔：有任务在下载时用较短的间隔，让进度看起来是「实时」的；
  // 空闲时放慢，减少无谓请求。
  // 这里用「每次请求完重新计算下一次延迟」的方式，
  // 而不是固定 setInterval，因此下载开始/结束能立刻改变节奏。
  const POLL_INTERVAL_ACTIVE = 2000
  const POLL_INTERVAL_IDLE = 5000

  function nextInterval() {
    return tasks.value.some(t => t.status === 'downloading' || t.status === 'pending')
      ? POLL_INTERVAL_ACTIVE
      : POLL_INTERVAL_IDLE
  }

  // 统一获取任务列表（始终带上当前筛选条件）
  async function fetchTasks() {
    loading.value = true
    try {
      const params = statusFilter.value
        ? `?status=${encodeURIComponent(statusFilter.value)}`
        : ''
      const newTasks = await get(`/tasks${params}`)
      // 只在数据变化时更新，避免不必要的重渲染
      if (JSON.stringify(tasks.value) !== JSON.stringify(newTasks)) {
        tasks.value = newTasks
      }
    } catch (e) {
      console.error('获取任务列表失败:', e)
    } finally {
      loading.value = false
    }
  }

  // 切换筛选条件并立即刷新
  function setStatusFilter(status) {
    statusFilter.value = status || ''
    return fetchTasks()
  }

  // 仅清空筛选条件（用于离开任务页时复位，不触发请求）
  function resetStatusFilter() {
    statusFilter.value = ''
  }

  function scheduleNext() {
    pollingTimer = setTimeout(async () => {
      if (!polling) return
      await fetchTasks()
      if (polling) scheduleNext()
    }, nextInterval())
  }

  // 启动轮询（单例模式）
  // 注意：即使已在轮询中也会先立即刷新一次，
  // 保证页面切换/筛选变化后立刻拿到最新数据
  function startPolling() {
    if (polling) {
      fetchTasks()
      return
    }
    polling = true
    fetchTasks()
    scheduleNext()
  }

  // 停止轮询
  function stopPolling() {
    polling = false
    if (pollingTimer) {
      clearTimeout(pollingTimer)
      pollingTimer = null
    }
  }

  async function createTask(args) {
    const res = await post('/tasks', args)
    // 有筛选条件时，新任务不一定属于当前列表，交给服务端刷新更准确
    if (statusFilter.value && res.status !== statusFilter.value) {
      fetchTasks()
    } else {
      tasks.value.unshift(res)
    }
    return res
  }

  async function deleteTask(id) {
    await del(`/tasks/${id}`)
    tasks.value = tasks.value.filter(t => t.id !== id)
  }

  // 取消正在进行的任务（记录保留，可再次重试）
  async function cancelTask(id) {
    const updated = await post(`/tasks/${id}/cancel`)
    replaceTask(updated)
    return updated
  }

  // 重试失败/中断的任务
  async function retryTask(id) {
    const updated = await post(`/tasks/${id}/retry`)
    replaceTask(updated)
    return updated
  }

  // 批量删除：本地先移除，再用服务端结果兜底刷新
  async function deleteTasks(ids) {
    const res = await del('/tasks', { ids })
    const deleted = new Set(res?.deleted || [])
    if (deleted.size > 0) {
      tasks.value = tasks.value.filter(t => !deleted.has(t.id))
    }
    if (res?.failed?.length) {
      fetchTasks()
    }
    return res
  }

  function replaceTask(updated) {
    if (!updated || updated.id == null) return
    const index = tasks.value.findIndex(t => t.id === updated.id)
    if (index !== -1) {
      tasks.value[index] = updated
    }
  }

  async function getTaskProgress(id) {
    return await get(`/tasks/${id}`)
  }

  async function getTaskLog(id) {
    return await get(`/tasks/${id}/log`)
  }

  return {
    tasks,
    loading,
    statusFilter,
    fetchTasks,
    setStatusFilter,
    resetStatusFilter,
    startPolling,
    stopPolling,
    createTask,
    deleteTask,
    cancelTask,
    retryTask,
    deleteTasks,
    getTaskProgress,
    getTaskLog
  }
})
