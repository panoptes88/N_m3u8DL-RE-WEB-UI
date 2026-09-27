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
  const pollingInterval = 5000 // 轮询间隔 5 秒

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

  // 启动轮询（单例模式）
  // 注意：即使已在轮询中也会先立即刷新一次，
  // 保证页面切换/筛选变化后立刻拿到最新数据
  function startPolling() {
    fetchTasks()
    if (pollingTimer) {
      return // 定时器已存在，不重复创建
    }
    pollingTimer = setInterval(() => {
      fetchTasks()
    }, pollingInterval)
  }

  // 停止轮询
  function stopPolling() {
    if (pollingTimer) {
      clearInterval(pollingTimer)
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
    getTaskProgress,
    getTaskLog
  }
})
