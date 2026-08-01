<template>
  <span class="status-tag" :style="tagStyle">
    <span class="status-dot" :class="{ pulsing: status === 'downloading' }" />
    <span>{{ text }}</span>
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  status: {
    type: String,
    required: true
  }
})

const STATUS_META = {
  pending: { text: '等待中', color: '#f59e0b' },
  downloading: { text: '下载中', color: '#6366f1' },
  completed: { text: '已完成', color: '#10b981' },
  failed: { text: '下载失败', color: '#ef4444' },
  interrupted: { text: '已中断', color: '#f97316' }
}

const meta = computed(() => STATUS_META[props.status] || { text: props.status, color: '#98a2b8' })
const text = computed(() => meta.value.text)
const tagStyle = computed(() => ({
  color: meta.value.color,
  backgroundColor: `${meta.value.color}1f`,
  '--dot-color': meta.value.color
}))
</script>

<style scoped>
.status-tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 2px 10px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 500;
  line-height: 20px;
  white-space: nowrap;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--dot-color);
}

.status-dot.pulsing {
  animation: pulse 1.2s ease-in-out infinite;
}

@keyframes pulse {
  50% {
    opacity: 0.3;
  }
}
</style>
