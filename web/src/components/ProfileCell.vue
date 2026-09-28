<template>
  <div class="profile-cell">
    <div class="profile-line">
      <a-tooltip :title="name" placement="topLeft">
        <span class="profile-name" @click="$emit('copy', name)">{{ name }}</span>
      </a-tooltip>
      <CopyOutlined class="copy-icon" title="复制方案名称" @click="$emit('copy', name)" />
    </div>
    <!-- 域名与方案名称相同时不重复显示：后端默认用域名当方案名，否则同一串长文本会出现两遍 -->
    <a-tooltip v-if="showDomain" :title="domain" placement="topLeft">
      <span class="profile-domain" @click="$emit('copy', domain)">{{ domain }}</span>
    </a-tooltip>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { CopyOutlined } from '@ant-design/icons-vue'

const props = defineProps({
  name: {
    type: String,
    default: ''
  },
  domain: {
    type: String,
    default: ''
  }
})

defineEmits(['copy'])

const showDomain = computed(() => !!props.domain && props.domain !== props.name)
</script>

<style scoped>
.profile-cell {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.profile-line {
  display: flex;
  align-items: center;
  gap: 4px;
  min-width: 0;
}

/* 长域名（如 cwfuommwagyrfcg6.automationengineering.space）必须省略号截断，
   否则会把列撑宽、把「操作」列挤出可视区 */
.profile-name {
  min-width: 0;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-1);
  cursor: pointer;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.profile-name:hover {
  color: var(--brand);
}

.profile-domain {
  display: block;
  min-width: 0;
  font-size: 12px;
  color: var(--text-3);
  cursor: pointer;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.profile-domain:hover {
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

.profile-line:hover .copy-icon {
  opacity: 1;
}

.copy-icon:hover {
  color: var(--brand);
}

@media (hover: none) {
  .copy-icon {
    opacity: 0.6;
  }
}
</style>
