<template>
  <span class="status-tag" :style="tagStyle">
    <span class="status-dot" :class="{ pulsing: status === 'downloading' }" />
    <span>{{ text }}</span>
  </span>
</template>

<script setup>
import { computed } from 'vue'
import { statusMeta } from '../utils/task'

const props = defineProps({
  status: {
    type: String,
    required: true
  }
})

const meta = computed(() => statusMeta(props.status))
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
