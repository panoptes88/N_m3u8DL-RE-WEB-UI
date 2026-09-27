<template>
  <div class="task-cell">
    <div class="task-line">
      <a-tooltip :title="outputName" placement="topLeft">
        <span class="task-name" @click="$emit('copy', outputName)">{{ outputName || '(未命名)' }}</span>
      </a-tooltip>
      <CopyOutlined
        v-if="outputName"
        class="copy-icon"
        title="复制文件名"
        @click="$emit('copy', outputName)"
      />
    </div>
    <a-tooltip :title="url" placement="topLeft">
      <span class="task-url" @click="$emit('copy', url)">{{ url }}</span>
    </a-tooltip>
  </div>
</template>

<script setup>
import { CopyOutlined } from '@ant-design/icons-vue'

defineProps({
  outputName: {
    type: String,
    default: ''
  },
  url: {
    type: String,
    default: ''
  }
})

defineEmits(['copy'])
</script>

<style scoped>
/* 任务单元格：输出文件名为主、URL 为辅，避免两列分别被截断 */
.task-cell {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.task-line {
  display: flex;
  align-items: center;
  gap: 4px;
  min-width: 0;
}

.task-name {
  min-width: 0;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-1);
  cursor: pointer;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.task-name:hover {
  color: var(--brand);
}

.task-url {
  display: block;
  min-width: 0;
  font-size: 12px;
  color: var(--text-3);
  cursor: pointer;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.task-url:hover {
  color: var(--brand);
}

.copy-icon {
  flex: none;
  font-size: 12px;
  color: var(--text-3);
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.15s, color 0.15s;
}

.task-line:hover .copy-icon {
  opacity: 1;
}

.copy-icon:hover {
  color: var(--brand);
}

/* 触屏设备没有 hover，直接常显 */
@media (hover: none) {
  .copy-icon {
    opacity: 0.6;
  }
}
</style>
