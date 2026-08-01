<template>
  <div class="brand-logo" :class="{ collapsed }">
    <div class="logo-tile">
      <CloudDownloadOutlined />
    </div>
    <div v-if="!collapsed" class="logo-text">
      <span class="logo-name">N_m3u8DL-RE</span>
      <span class="logo-sub">Web UI</span>
    </div>
  </div>
</template>

<script setup>
import { CloudDownloadOutlined } from '@ant-design/icons-vue'

defineProps({
  collapsed: {
    type: Boolean,
    default: false
  }
})
</script>

<style scoped>
.brand-logo {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 16px;
  height: 60px;
  overflow: hidden;
}

.brand-logo.collapsed {
  justify-content: center;
  padding: 0;
}

.logo-tile {
  flex-shrink: 0;
  width: 34px;
  height: 34px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  color: #fff;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  box-shadow: 0 4px 10px rgba(99, 102, 241, 0.35);
}

.logo-text {
  display: flex;
  flex-direction: column;
  line-height: 1.25;
  white-space: nowrap;
}

.logo-name {
  font-size: 15px;
  font-weight: 700;
  color: var(--text-1);
  letter-spacing: 0.2px;
}

.logo-sub {
  font-size: 11px;
  color: var(--text-3);
}
</style>
