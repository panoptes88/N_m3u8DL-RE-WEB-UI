<template>
  <a-modal
    :open="open"
    title="任务日志"
    :footer="null"
    :width="860"
    @update:open="$emit('update:open', $event)"
  >
    <div class="log-toolbar">
      <a-space :size="8">
        <span v-if="isLive" class="live-badge">
          <span class="live-dot" />
          实时刷新中
        </span>
        <span v-else class="log-hint">任务已结束</span>
        <span class="log-hint">共 {{ lineCount }} 行</span>
      </a-space>
      <a-space :size="8">
        <a-checkbox v-model:checked="autoScroll">自动滚动</a-checkbox>
        <a-button size="small" :disabled="!content" @click="copyAll">
          <template #icon><CopyOutlined /></template>
          复制
        </a-button>
        <a-button size="small" :disabled="!content" @click="downloadLog">
          <template #icon><DownloadOutlined /></template>
          下载
        </a-button>
      </a-space>
    </div>

    <a-spin :spinning="loading">
      <pre ref="logRef" class="log-content" @scroll="handleScroll">{{ content || '暂无日志' }}</pre>
    </a-spin>

    <transition name="fade">
      <a-button
        v-if="!autoScroll && content"
        class="back-to-bottom"
        size="small"
        type="primary"
        @click="scrollToBottom"
      >
        回到底部
      </a-button>
    </transition>
  </a-modal>
</template>

<script setup>
import { ref, computed, watch, nextTick, onUnmounted } from 'vue'
import { message } from 'ant-design-vue'
import { CopyOutlined, DownloadOutlined } from '@ant-design/icons-vue'
import { useTaskStore } from '../stores/task'
import { copyToClipboard } from '../utils/clipboard'

const props = defineProps({
  open: {
    type: Boolean,
    default: false
  },
  taskId: {
    type: [Number, String],
    default: null
  },
  // 任务状态：downloading/pending 时持续刷新
  status: {
    type: String,
    default: ''
  }
})

defineEmits(['update:open'])

const taskStore = useTaskStore()
const loading = ref(false)
const content = ref('')
const logRef = ref(null)
const autoScroll = ref(true)

// 请求序号：防止切换任务后，较慢的旧响应覆盖当前任务的日志
let requestSeq = 0
let timer = null
const REFRESH_INTERVAL = 2000

const isLive = computed(() => props.status === 'downloading' || props.status === 'pending')
const lineCount = computed(() => (content.value ? content.value.split('\n').length : 0))

async function fetchLog({ silent = false } = {}) {
  if (props.taskId == null) return
  const seq = ++requestSeq
  if (!silent) loading.value = true
  try {
    const res = await taskStore.getTaskLog(props.taskId)
    if (seq !== requestSeq || !props.open) return
    const next = res.log || ''
    // 内容未变化时不触碰 DOM，避免打断用户的滚动位置
    if (next !== content.value) {
      content.value = next
      if (autoScroll.value) scrollToBottom()
    }
  } catch {
    if (seq === requestSeq && !silent) message.error('获取日志失败')
  } finally {
    if (seq === requestSeq && !silent) loading.value = false
  }
}

function stopTimer() {
  if (timer) {
    clearInterval(timer)
    timer = null
  }
}

function startTimer() {
  stopTimer()
  timer = setInterval(() => fetchLog({ silent: true }), REFRESH_INTERVAL)
}

async function scrollToBottom() {
  await nextTick()
  const el = logRef.value
  if (el) el.scrollTop = el.scrollHeight
}

// 用户手动上滚时暂停自动滚动，避免正在查看历史日志时被拽到底部
function handleScroll() {
  const el = logRef.value
  if (!el) return
  const atBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 24
  if (atBottom) {
    autoScroll.value = true
  } else if (autoScroll.value) {
    autoScroll.value = false
  }
}

async function copyAll() {
  await copyToClipboard(content.value)
}

function downloadLog() {
  const blob = new Blob([content.value], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `task_${props.taskId}.log`
  link.click()
  URL.revokeObjectURL(url)
}

watch(
  () => [props.open, props.taskId, props.status],
  async ([open]) => {
    if (!open || props.taskId == null) {
      stopTimer()
      return
    }
    autoScroll.value = true
    await fetchLog()
    // 仅在任务进行中才轮询，结束后自动停止
    if (isLive.value) startTimer()
    else stopTimer()
  },
  { immediate: true }
)

onUnmounted(() => {
  stopTimer()
  requestSeq++
})
</script>

<style scoped>
.log-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 10px;
}

.log-hint {
  font-size: 12px;
  color: var(--text-2);
}

.live-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 1px 8px;
  border-radius: 999px;
  font-size: 12px;
  color: var(--brand);
  background: rgba(99, 102, 241, 0.1);
}

.live-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--brand);
  animation: live-pulse 1.2s ease-in-out infinite;
}

@keyframes live-pulse {
  50% {
    opacity: 0.25;
  }
}

.log-content {
  margin: 0;
  padding: 12px 14px;
  max-height: 60vh;
  overflow: auto;
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
  font-size: 12px;
  line-height: 1.7;
  white-space: pre-wrap;
  word-break: break-all;
  border-radius: var(--radius-sm);
  background: rgba(100, 116, 139, 0.08);
  color: var(--text-1);
}

:global(.dark) .log-content {
  background: rgba(255, 255, 255, 0.04);
}

.back-to-bottom {
  position: absolute;
  left: 50%;
  bottom: 22px;
  transform: translateX(-50%);
  box-shadow: var(--card-shadow-hover);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
