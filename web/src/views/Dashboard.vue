<template>
  <div class="dashboard">
    <PageHeader title="首页" subtitle="下载任务概览与快速创建" />

    <a-row :gutter="[16, 16]" class="stats-row">
      <a-col v-for="stat in stats" :key="stat.label" :xs="24" :sm="8">
        <div class="stat-card">
          <div class="stat-icon" :style="{ background: `${stat.color}1a`, color: stat.color }">
            <component :is="stat.icon" />
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stat.value }}</div>
            <div class="stat-label">{{ stat.label }}</div>
          </div>
        </div>
      </a-col>
    </a-row>

    <a-card title="快速下载" class="app-card quick-download">
      <a-form layout="vertical" :model="quickDownloadForm">
        <a-form-item
          label="m3u8 URL"
          :rules="[{ required: true, message: '请输入m3u8链接' }]"
        >
          <a-input
            v-model:value="quickDownloadForm.url"
            placeholder="https://example.com/video.m3u8"
            size="large"
          />
        </a-form-item>
        <a-row :gutter="12">
          <a-col :xs="24" :sm="12">
            <a-form-item label="输出名称（可选）">
              <a-input
                v-model:value="quickDownloadForm.outputName"
                placeholder="output.mp4"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12" class="btn-col">
            <a-form-item label=" ">
              <a-button
                type="primary"
                size="large"
                block
                :loading="quickDownloadLoading"
                @click="handleQuickDownload"
              >
                <template #icon><CloudDownloadOutlined /></template>
                开始下载
              </a-button>
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-card>

    <a-card title="最近任务" class="app-card recent-tasks">
      <template #extra>
        <a-button type="link" @click="$router.push('/tasks')">查看全部</a-button>
      </template>
      <a-table
        :columns="columns"
        :data-source="recentTasks"
        :pagination="false"
        size="middle"
        :loading="loading"
        :scroll="{ x: 500 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <StatusTag :status="record.status" />
          </template>
          <template v-if="column.key === 'progress'">
            <a-progress
              :percent="record.progress"
              :status="record.status === 'failed' ? 'exception' : 'active'"
              :stroke-color="record.status === 'failed' ? undefined : progressGradient"
              size="small"
            />
          </template>
          <template v-if="column.key === 'action'">
            <a-space :size="4">
              <a-button
                size="small"
                type="text"
                @click="viewLog(record)"
              >
                日志
              </a-button>
              <a-popconfirm
                title="确定删除此任务？"
                @confirm="deleteTask(record.id)"
              >
                <a-button size="small" type="text" danger>删除</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <TaskLogModal v-model:open="logModalVisible" :task-id="logTaskId" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import {
  ClockCircleOutlined,
  CloudDownloadOutlined,
  CheckCircleOutlined
} from '@ant-design/icons-vue'
import { useTaskStore } from '../stores/task'
import PageHeader from '../components/PageHeader.vue'
import StatusTag from '../components/StatusTag.vue'
import TaskLogModal from '../components/TaskLogModal.vue'

const taskStore = useTaskStore()

const loading = ref(false)
const quickDownloadLoading = ref(false)
const quickDownloadForm = ref({ url: '', outputName: '' })
const logModalVisible = ref(false)
const logTaskId = ref(null)

const progressGradient = { from: '#6366f1', to: '#8b5cf6' }

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
  { title: 'URL', dataIndex: 'url', key: 'url', ellipsis: true },
  { title: '状态', dataIndex: 'status', key: 'status', width: 90 },
  { title: '进度', dataIndex: 'progress', key: 'progress', width: 130 },
  { title: '操作', key: 'action', width: 120 }
]

const recentTasks = computed(() => taskStore.tasks.slice(0, 5))
const pendingCount = computed(() => taskStore.tasks.filter(t => t.status === 'pending').length)
const downloadingCount = computed(() => taskStore.tasks.filter(t => t.status === 'downloading').length)
const completedCount = computed(() => taskStore.tasks.filter(t => t.status === 'completed').length)

const stats = computed(() => [
  { label: '等待中任务', value: pendingCount.value, color: '#f59e0b', icon: ClockCircleOutlined },
  { label: '下载中任务', value: downloadingCount.value, color: '#6366f1', icon: CloudDownloadOutlined },
  { label: '已完成任务', value: completedCount.value, color: '#10b981', icon: CheckCircleOutlined }
])

async function handleQuickDownload() {
  if (!quickDownloadForm.value.url) {
    message.error('请输入m3u8链接')
    return
  }

  quickDownloadLoading.value = true
  try {
    await taskStore.createTask({
      url: quickDownloadForm.value.url,
      output_name: quickDownloadForm.value.outputName || ''
    })
    message.success('任务已创建')
    quickDownloadForm.value = { url: '', outputName: '' }
  } catch (err) {
    message.error(err.response?.data?.error || '创建任务失败')
  } finally {
    quickDownloadLoading.value = false
  }
}

function viewLog(task) {
  logTaskId.value = task.id
  logModalVisible.value = true
}

async function deleteTask(id) {
  try {
    await taskStore.deleteTask(id)
    message.success('删除成功')
  } catch {
    message.error('删除失败')
  }
}

onMounted(() => {
  // 使用 store 统一管理的轮询（单例模式）
  taskStore.startPolling()
})
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* 统计卡片 */
.stat-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 20px;
  border-radius: 12px;
  background: var(--surface);
  border: 1px solid var(--border-soft);
  box-shadow: var(--card-shadow);
  transition: transform 0.2s, box-shadow 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--card-shadow-hover);
}

.stat-icon {
  flex-shrink: 0;
  width: 46px;
  height: 46px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
}

.stat-value {
  font-size: 26px;
  font-weight: 700;
  line-height: 1.2;
  color: var(--text-1);
  font-variant-numeric: tabular-nums;
}

.stat-label {
  margin-top: 2px;
  font-size: 13px;
  color: var(--text-2);
}

.btn-col {
  display: flex;
  align-items: flex-end;
}

.recent-tasks {
  flex: 1;
}

@media (max-width: 576px) {
  .stat-card {
    padding: 16px;
  }

  .stat-value {
    font-size: 22px;
  }

  .recent-tasks :deep(.ant-table-cell) {
    padding: 8px !important;
  }
}
</style>
