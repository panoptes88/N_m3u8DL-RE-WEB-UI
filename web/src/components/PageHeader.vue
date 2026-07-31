<template>
  <div class="page-header">
    <div class="page-header-text">
      <h2 class="page-title">{{ title }}</h2>
      <p v-if="subtitle" class="page-subtitle">{{ subtitle }}</p>
    </div>
    <div v-if="$slots.default" class="page-header-extra">
      <slot />
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true
  },
  subtitle: {
    type: String,
    default: ''
  }
})
</script>

<style scoped>
.page-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 20px;
}

.page-title {
  margin: 0;
  font-size: 20px;
  font-weight: 700;
  color: var(--text-1);
  letter-spacing: 0.2px;
}

.page-subtitle {
  margin: 4px 0 0;
  font-size: 13px;
  color: var(--text-2);
}

.page-header-extra {
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>
