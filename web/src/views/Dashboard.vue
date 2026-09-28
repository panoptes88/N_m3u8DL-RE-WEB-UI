<template>
  <div class="dashboard">
    <PageHeader title="首页" subtitle="下载任务概览与快速创建" />

    <a-row :gutter="[16, 16]" class="stats-row">
      <a-col v-for="stat in stats" :key="stat.label" :xs="12" :sm="12" :md="6">
        <div
          class="stat-card"
          role="button"
          tabindex="0"
          :title="`查看${stat.label}的任务`"
          @click="goTasks(stat.filter)"
          @keydown.enter="goTasks(stat.filter)"
          @keydown.space.prevent="goTasks(stat.filter)"
        >
          <a-skeleton
            v-if="initialLoading"
            active
            :title="false"
            :paragraph="{ rows: 2, width: ['38%', '68%'] }"
          />
          <template v-else>
            <div class="stat-icon" :style="{ background: `${stat.color}1a`, color: stat.color }">
              <component :is="stat.icon" />
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stat.value }}</div>
              <div class="stat-label">{{ stat.label }}</div>
            </div>
          </template>
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
            <a-form-item label="输出名称（可选，留空自动生成）">
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

    <a-card class="app-card recent-tasks">
      <template #title>
        最近任务
        <span class="card-title-hint">共 {{ taskStore.tasks.length }} 个</span>
      </template>
      <template #extra>
        <a-button type="link" @click="$router.push('/tasks')">查看全部</a-button>
      </template>
      <a-table
        :columns="columns"
        :data-source="recentTasks"
        :pagination="false"
        size="middle"
        :loading="taskStore.loading"
        :scroll="{ x: 760 }"
        :row-class-name="rowClassName"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'task'">
            <TaskCell
              :output-name="record.output_name"
              :url="record.url"
              @copy="copyToClipboard"
            />
          </template>
          <template v-if="column.key === 'status'">
            <a-space :size="4">
              <StatusTag :status="record.status" />
              <a-tooltip v-if="record.error_msg" :title="record.error_msg" placement="topLeft">
                <ExclamationCircleOutlined class="error-hint" />
              </a-tooltip>
            </a-space>
          </template>
          <template v-if="column.key === 'progress'">
            <div class="progress-cell">
              <a-progress
                :percent="record.progress"
                :status="isErrorStatus(record.status) ? 'exception' : 'active'"
                :stroke-color="isErrorStatus(record.status) ? undefined : progressGradient"
                :show-info="false"
                size="small"
              />
              <div class="progress-meta">
                <span class="progress-percent">{{ record.progress }}%</span>
                <span v-if="progressMeta(record)" class="progress-detail">{{ progressMeta(record) }}</span>
              </div>
            </div>
          </template>
          <template v-if="column.key === 'duration'">
            {{ formatDuration(record) }}
          </template>
          <template v-if="column.key === 'action'">
            <TaskActions
              :task="record"
              @log="viewLog"
              @cancel="handleCancel"
              @retry="handleRetry"
              @delete="confirmDelete"
            />
          </template>
        </template>
        <template #emptyText>
          <a-empty description="还没有任务，在上方填入链接即可开始下载" />
        </template>
      </a-table>
    </a-card>

    <TaskLogModal
      v-model:open="logModalVisible"
      :task-id="logTaskId"
      :status="logTaskStatus"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import {
  ClockCircleOutlined,
  CloudDownloadOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  ExclamationCircleOutlined
} from '@ant-design/icons-vue'
import { useTaskStore } from '../stores/task'
import { formatDuration, progressMeta, isErrorStatus } from '../utils/task'
import { copyToClipboard } from '../utils/clipboard'
import PageHeader from '../components/PageHeader.vue'
import StatusTag from '../components/StatusTag.vue'
import TaskActions from '../components/TaskActions.vue'
import TaskCell from '../components/TaskCell.vue'
import TaskLogModal from '../components/TaskLogModal.vue'

const router = useRouter()
const taskStore = useTaskStore()

const quickDownloadLoading = ref(false)
const quickDownloadForm = ref({ url: '', outputName: '' })
const logModalVisible = ref(false)
const logTaskId = ref(null)

// 打开中的日志对应的任务状态：随轮询更新，用于决定是否持续刷新日志
const logTaskStatus = computed(
  () => taskStore.tasks.find(t => t.id === logTaskId.value)?.status || ''
)

const progressGradient = { from: '#6366f1', to: '#8b5cf6' }

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
  { title: '任务', key: 'task', width: 300 },
  { title: '状态', dataIndex: 'status', key: 'status', width: 130 },
  { title: '进度', dataIndex: 'progress', key: 'progress', width: 200 },
  { title: '耗时', key: 'duration', width: 100, responsive: ['lg'] },
  { title: '操作', key: 'action', width: 150 }
]

const recentTasks = computed(() => taskStore.tasks.slice(0, 5))

// 首次加载（还没有任何数据）时用骨架屏代替「0」的闪烁
const initialLoading = computed(() => taskStore.loading && taskStore.tasks.length === 0)

// 统计口径覆盖全部状态，避免出现「各卡片之和 ≠ 任务总数」
const countByStatus = statuses =>
  taskStore.tasks.filter(t => statuses.includes(t.status)).length

const stats = computed(() => [
  {
    label: '下载中',
    value: countByStatus(['downloading']),
    color: '#6366f1',
    icon: CloudDownloadOutlined,
    filter: 'downloading'
  },
  {
    label: '等待中',
    value: countByStatus(['pending']),
    color: '#f59e0b',
    icon: ClockCircleOutlined,
    filter: 'pending'
  },
  {
    label: '已完成',
    value: countByStatus(['completed']),
    color: '#10b981',
    icon: CheckCircleOutlined,
    filter: 'completed'
  },
  {
    label: '失败 / 中断',
    value: countByStatus(['failed', 'interrupted']),
    color: '#ef4444',
    icon: CloseCircleOutlined,
    filter: 'failed,interrupted'
  }
])

function goTasks(status) {
  router.push({ name: 'Tasks', query: status ? { status } : {} })
}

// 异常结束的任务整行淡色区分，便于快速扫视
function rowClassName(record) {
  if (record.status === 'failed') return 'row-failed'
  if (record.status === 'interrupted') return 'row-interrupted'
  return ''
}

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

// 删除不可恢复，二次确认（入口在下拉菜单里）
function confirmDelete(task) {
  Modal.confirm({
    title: '确定删除此任务？',
    content: `「${task.output_name || task.url}」将被删除，此操作不可恢复。`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    onOk: () => deleteTask(task.id)
  })
}

async function handleCancel(task) {
  try {
    await taskStore.cancelTask(task.id)
    message.success('已取消任务')
  } catch (err) {
    message.error(err.response?.data?.error || '取消失败')
  }
}

async function handleRetry(task) {
  try {
    await taskStore.retryTask(task.id)
    message.success('已加入下载队列')
  } catch (err) {
    message.error(err.response?.data?.error || '重试失败')
  }
}

onMounted(() => {
  // 使用 store 统一管理的轮询（单例模式）
  taskStore.startPolling()
})

onUnmounted(() => {
  // 离开首页即停止轮询，避免在登录页等场景继续后台请求
  taskStore.stopPolling()
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
  border-radius: var(--radius-md);
  background: var(--surface);
  border: 1px solid var(--border-soft);
  box-shadow: var(--card-shadow);
  transition: transform 0.2s, box-shadow 0.2s;
  cursor: pointer;
}

.stat-card:hover,
.stat-card:focus-visible {
  transform: translateY(-2px);
  box-shadow: var(--card-shadow-hover);
  border-color: rgba(99, 102, 241, 0.35);
  outline: none;
}

.stat-icon {
  flex-shrink: 0;
  width: 46px;
  height: 46px;
  border-radius: var(--radius-md);
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

.card-title-hint {
  margin-left: 8px;
  font-size: 12px;
  font-weight: 400;
  color: var(--text-2);
}

/* 进度单元格：细进度条 + 下方速度/大小 */
.progress-cell {
  min-width: 0;
}

.progress-cell :deep(.ant-progress) {
  width: 100%;
  margin-bottom: 0;
  line-height: 1;
}

.progress-meta {
  display: flex;
  align-items: baseline;
  gap: 8px;
  margin-top: 4px;
  font-size: 12px;
  font-variant-numeric: tabular-nums;
}

.progress-percent {
  font-weight: 600;
  color: var(--text-1);
}

.progress-detail {
  color: var(--text-2);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.error-hint {
  color: #ef4444;
  cursor: help;
}

.recent-tasks :deep(.row-failed) > td {
  background: rgba(239, 68, 68, 0.05);
}

.recent-tasks :deep(.row-interrupted) > td {
  background: rgba(249, 115, 22, 0.06);
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
