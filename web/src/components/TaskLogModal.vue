<template>
  <a-modal
    :open="open"
    title="任务日志"
    :footer="null"
    :width="760"
    @update:open="$emit('update:open', $event)"
  >
    <a-spin :spinning="loading">
      <pre class="log-content">{{ content || '暂无日志' }}</pre>
    </a-spin>
  </a-modal>
</template>

<script setup>
import { ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import { useTaskStore } from '../stores/task'

const props = defineProps({
  open: {
    type: Boolean,
    default: false
  },
  taskId: {
    type: [Number, String],
    default: null
  }
})

defineEmits(['update:open'])

const taskStore = useTaskStore()
const loading = ref(false)
const content = ref('')

// 请求序号：防止切换任务后，较慢的旧响应覆盖当前任务的日志
let requestSeq = 0

watch(
  () => [props.open, props.taskId],
  async ([open, taskId]) => {
    if (!open || taskId == null) return
    const seq = ++requestSeq
    loading.value = true
    content.value = ''
    try {
      const res = await taskStore.getTaskLog(taskId)
      if (seq !== requestSeq || !props.open) return
      content.value = res.log || ''
    } catch {
      if (seq !== requestSeq) return
      message.error('获取日志失败')
    } finally {
      if (seq === requestSeq) {
        loading.value = false
      }
    }
  }
)
</script>

<style scoped>
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
  border-radius: 8px;
  background: rgba(100, 116, 139, 0.08);
  color: var(--text-1);
}

:global(.dark) .log-content {
  background: rgba(255, 255, 255, 0.04);
}
</style>
