<template>
  <div class="tasks-page">
    <PageHeader title="下载任务" subtitle="创建任务并跟踪下载进度" />

    <!-- 创建任务表单 -->
    <a-card title="创建下载任务" class="app-card">
      <template #extra>
        <a-space>
          <a-select
            v-model:value="selectedProfileId"
            placeholder="选择下载方案"
            style="width: 200px"
            allow-clear
            @change="handleProfileChange"
          >
            <a-select-option v-for="profile in profileStore.profiles" :key="profile.id" :value="profile.id">
              {{ profile.name }}
            </a-select-option>
          </a-select>
          <a-button @click="showProfileManager = true">
            <template #icon><SettingOutlined /></template>
            管理方案
          </a-button>
        </a-space>
      </template>
      <a-form
        ref="formRef"
        :model="formState"
        :rules="formRules"
        layout="vertical"
      >
        <a-row :gutter="16">
          <a-col :span="24">
            <a-form-item label="m3u8 URL" name="url">
              <a-input
                v-model:value="formState.url"
                placeholder="https://example.com/video.m3u8"
                size="large"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="16">
          <a-col :xs="24" :sm="12" :md="8">
            <a-form-item label="输出文件名" name="outputName">
              <a-input
                v-model:value="formState.outputName"
                placeholder="output.mp4"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12" :md="8">
            <a-form-item label="线程数" name="threadCount">
              <a-input-number v-model:value="formState.threadCount" :min="1" :max="128" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12" :md="8">
            <a-form-item label="重试次数" name="retryCount">
              <a-input-number v-model:value="formState.retryCount" :min="0" :max="100" style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="16">
          <a-col :xs="24" :sm="12">
            <a-form-item label="请求头" name="headers">
              <a-input
                v-model:value="formState.headers"
                placeholder='如: Cookie: xxx; User-Agent: xxx'
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12">
            <a-form-item label="Base URL" name="baseUrl">
              <a-input
                v-model:value="formState.baseUrl"
                placeholder="可选，用于补全相对路径"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="16">
          <a-col :xs="24" :sm="12" :md="8">
            <a-form-item label=" " name="delAfterDone">
              <a-checkbox v-model:checked="formState.delAfterDone">
                下载完成后删除临时文件
              </a-checkbox>
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12" :md="8">
            <a-form-item label=" " name="binaryMerge">
              <a-checkbox v-model:checked="formState.binaryMerge">
                启用二进制合并
              </a-checkbox>
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12" :md="8">
            <a-form-item label=" " name="autoSelect">
              <a-checkbox v-model:checked="formState.autoSelect">
                自动选择最佳轨道
              </a-checkbox>
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="16">
          <a-col :xs="24" :sm="12" :md="8">
            <a-form-item label=" " name="skipSegmentsCheck">
              <a-checkbox v-model:checked="formState.skipSegmentsCheck">
                跳过完整性检测
              </a-checkbox>
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12" :md="8">
            <a-form-item label=" " name="concurrentDownload">
              <a-checkbox v-model:checked="formState.concurrentDownload">
                并行下载音视频
              </a-checkbox>
            </a-form-item>
          </a-col>
        </a-row>

        <!-- 解密选项 -->
        <div class="form-section-title">解密选项</div>
        <a-row :gutter="16">
          <a-col :xs="24" :sm="12" :md="8">
            <a-form-item label="解密密钥" name="key">
              <a-input
                v-model:value="formState.key"
                placeholder="KID:KEY 或直接 KEY"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12" :md="8">
            <a-form-item label="解密引擎" name="decryptionEngine">
              <a-select v-model:value="formState.decryptionEngine">
                <a-select-option value="MP4DECRYPT">MP4DECRYPT</a-select-option>
                <a-select-option value="FFMPEG">FFMPEG</a-select-option>
                <a-select-option value="SHAKA_PACKAGER">SHAKA_PACKAGER</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <!-- 代理设置 -->
        <div class="form-section-title">代理设置</div>
        <a-row :gutter="16">
          <a-col :xs="24" :sm="24" :md="12">
            <a-form-item label="自定义代理" name="customProxy">
              <a-input
                v-model:value="formState.customProxy"
                placeholder="如: http://127.0.0.1:7890"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <!-- 自定义参数 -->
        <div class="form-section-title">其他参数</div>
        <a-row :gutter="16">
          <a-col :span="24">
            <a-form-item label="自定义参数" name="customArgs">
              <a-input
                v-model:value="formState.customArgs"
                placeholder="其他命令行参数，如: --log-level DEBUG"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item style="margin-bottom: 0">
          <a-space>
            <a-button type="primary" size="large" :loading="creating" @click="handleCreate">
              创建任务
            </a-button>
            <a-button size="large" @click="resetForm">
              重置
            </a-button>
            <a-checkbox v-model:checked="keepFormAfterCreate">
              创建后保留表单
            </a-checkbox>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>

    <!-- 任务列表 -->
    <a-card title="下载任务" class="app-card task-list">
      <template #extra>
        <a-space>
          <a-select
            v-model:value="statusFilter"
            style="width: 120px"
            @change="fetchTasks"
          >
            <a-select-option value="">全部</a-select-option>
            <a-select-option value="pending">等待中</a-select-option>
            <a-select-option value="downloading">下载中</a-select-option>
            <a-select-option value="completed">已完成</a-select-option>
            <a-select-option value="failed">失败</a-select-option>
          </a-select>
          <a-button @click="fetchTasks">
            <template #icon><ReloadOutlined /></template>
            刷新
          </a-button>
        </a-space>
      </template>

      <a-table
        :columns="columns"
        :data-source="taskStore.tasks"
        :pagination="{ pageSize: 10 }"
        :loading="taskStore.loading"
        :scroll="{ x: 800 }"
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
            />
          </template>
          <template v-if="column.key === 'createdAt'">
            {{ formatTime(record.created_at) }}
          </template>
          <template v-if="column.key === 'action'">
            <a-space :size="4">
              <a-button size="small" type="text" @click="viewLog(record)">日志</a-button>
              <a-button size="small" type="text" @click="saveAsProfile(record)">保存方案</a-button>
              <a-popconfirm
                title="确定删除此任务？"
                @confirm="handleDelete(record.id)"
              >
                <a-button size="small" type="text" danger>删除</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 日志弹窗 -->
    <TaskLogModal v-model:open="logModalVisible" :task-id="logTaskId" />

    <!-- 方案管理弹窗 -->
    <a-modal
      v-model:open="showProfileManager"
      title="下载方案管理"
      :footer="null"
      width="800px"
    >
      <a-table
        :columns="profileColumns"
        :data-source="profileStore.profiles"
        :pagination="{ pageSize: 10 }"
        :loading="profileStore.loading"
        :scroll="{ x: 600 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <div v-if="editingProfileId === record.id">
              <a-input
                v-model:value="editingProfileName"
                size="small"
                @blur="saveProfileName(record)"
                @pressEnter="$event.target.blur()"
                autofocus
              />
            </div>
            <div v-else @dblclick="startEditProfileName(record)" style="cursor: pointer">
              {{ record.name }}
            </div>
          </template>
          <template v-if="column.key === 'created_at'">
            {{ formatDate(record.created_at) }}
          </template>
          <template v-if="column.key === 'action'">
            <a-space :size="4">
              <a-button size="small" type="text" @click="loadProfileToForm(record)">加载</a-button>
              <a-button size="small" type="text" @click="viewProfileDetail(record)">详情</a-button>
              <a-popconfirm
                title="确定删除此方案？"
                @confirm="handleDeleteProfile(record.id)"
              >
                <a-button size="small" type="text" danger>删除</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-modal>

    <!-- 方案详情弹窗 -->
    <a-modal
      v-model:open="showProfileDetail"
      title="方案详情"
      :footer="null"
      width="600px"
    >
      <a-descriptions v-if="currentProfile" bordered :column="1">
        <a-descriptions-item label="ID">{{ currentProfile.id }}</a-descriptions-item>
        <a-descriptions-item label="方案名称">{{ currentProfile.name }}</a-descriptions-item>
        <a-descriptions-item label="域名">{{ currentProfile.domain || '-' }}</a-descriptions-item>
        <a-descriptions-item label="线程数">{{ currentProfile.thread_count }}</a-descriptions-item>
        <a-descriptions-item label="重试次数">{{ currentProfile.retry_count }}</a-descriptions-item>
        <a-descriptions-item label="请求头">{{ currentProfile.headers || '-' }}</a-descriptions-item>
        <a-descriptions-item label="Base URL">{{ currentProfile.base_url || '-' }}</a-descriptions-item>
        <a-descriptions-item label="下载后删除临时文件">{{ currentProfile.del_after_done ? '是' : '否' }}</a-descriptions-item>
        <a-descriptions-item label="二进制合并">{{ currentProfile.binary_merge ? '是' : '否' }}</a-descriptions-item>
        <a-descriptions-item label="自动选择轨道">{{ currentProfile.auto_select ? '是' : '否' }}</a-descriptions-item>
        <a-descriptions-item label="跳过完整性检测">{{ currentProfile.skip_segments_check ? '是' : '否' }}</a-descriptions-item>
        <a-descriptions-item label="并行下载音视频">{{ currentProfile.concurrent_download ? '是' : '否' }}</a-descriptions-item>
        <a-descriptions-item label="解密引擎">{{ currentProfile.decryption_engine }}</a-descriptions-item>
        <a-descriptions-item label="自定义参数">{{ currentProfile.custom_args || '-' }}</a-descriptions-item>
        <a-descriptions-item label="自定义代理">{{ currentProfile.custom_proxy || '-' }}</a-descriptions-item>
        <a-descriptions-item label="创建时间">{{ formatDate(currentProfile.created_at) }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { ReloadOutlined, SettingOutlined } from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'
import { useTaskStore } from '../stores/task'
import { useProfileStore } from '../stores/profile'
import PageHeader from '../components/PageHeader.vue'
import StatusTag from '../components/StatusTag.vue'
import TaskLogModal from '../components/TaskLogModal.vue'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const taskStore = useTaskStore()
const profileStore = useProfileStore()

const creating = ref(false)
const formRef = ref(null)
const statusFilter = ref('')
const keepFormAfterCreate = ref(false) // 创建后保留表单
const logModalVisible = ref(false)
const logTaskId = ref(null)
const selectedProfileId = ref(null)
const showProfileManager = ref(false)
const showProfileDetail = ref(false)
const currentProfile = ref(null)
const editingProfileId = ref(null)
const editingProfileName = ref('')

const progressGradient = { from: '#6366f1', to: '#8b5cf6' }

// 表单数据
const formState = reactive({
  url: '',
  outputName: '',
  threadCount: 32,
  retryCount: 15,
  headers: '',
  baseUrl: '',
  delAfterDone: true,
  binaryMerge: false,
  autoSelect: false,
  skipSegmentsCheck: false,
  concurrentDownload: false,
  key: '',
  decryptionEngine: 'MP4DECRYPT',
  customArgs: '',
  customProxy: ''
})

const formRules = {
  url: [{ required: true, message: '请输入m3u8链接' }],
  outputName: [{ required: true, message: '请输入输出文件名' }]
}

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
  { title: 'URL', dataIndex: 'url', key: 'url', ellipsis: true, width: 200 },
  { title: '输出', dataIndex: 'output_name', key: 'outputName', ellipsis: true, width: 100 },
  { title: '状态', dataIndex: 'status', key: 'status', width: 90 },
  { title: '进度', dataIndex: 'progress', key: 'progress', width: 120 },
  { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', width: 150, responsive: ['lg'] },
  { title: '操作', key: 'action', width: 170 }
]

const profileColumns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 50 },
  { title: '方案名称', dataIndex: 'name', key: 'name', width: 150 },
  { title: '域名', dataIndex: 'domain', key: 'domain', width: 150 },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 100 },
  { title: '操作', key: 'action', width: 150 }
]

function formatTime(time) {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

function formatDate(time) {
  return dayjs(time).format('YYYY-MM-DD')
}

function viewProfileDetail(profile) {
  currentProfile.value = profile
  showProfileDetail.value = true
}

function startEditProfileName(profile) {
  editingProfileId.value = profile.id
  editingProfileName.value = profile.name
}

async function saveProfileName(profile) {
  if (!editingProfileName.value.trim()) {
    message.error('方案名称不能为空')
    editingProfileId.value = null
    return
  }

  if (editingProfileName.value !== profile.name) {
    try {
      await profileStore.updateProfile(profile.id, { name: editingProfileName.value })
      message.success('方案名称已更新')
    } catch {
      message.error('更新失败')
    }
  }
  editingProfileId.value = null
}

async function fetchTasks() {
  const status = statusFilter.value ? `?status=${statusFilter.value}` : ''
  await taskStore.fetchTasks(status)
}

function resetForm() {
  formState.url = ''
  formState.outputName = ''
  formState.threadCount = 32
  formState.retryCount = 15
  formState.headers = ''
  formState.baseUrl = ''
  formState.delAfterDone = true
  formState.binaryMerge = false
  formState.autoSelect = false
  formState.skipSegmentsCheck = false
  formState.concurrentDownload = false
  formState.key = ''
  formState.decryptionEngine = 'MP4DECRYPT'
  formState.customArgs = ''
  formState.customProxy = ''
  formRef.value?.clearValidate()
}

async function handleCreate() {
  try {
    await formRef.value.validate()
  } catch {
    return
  }

  creating.value = true
  try {
    // 构建参数字符串
    const args = {
      url: formState.url,
      output_name: formState.outputName,
      thread_count: formState.threadCount,
      retry_count: formState.retryCount,
      headers: formState.headers,
      base_url: formState.baseUrl,
      del_after_done: formState.delAfterDone,
      binary_merge: formState.binaryMerge,
      auto_select: formState.autoSelect,
      skip_segments_check: formState.skipSegmentsCheck,
      concurrent_download: formState.concurrentDownload,
      key: formState.key,
      decryption_engine: formState.decryptionEngine,
      custom_args: formState.customArgs,
      custom_proxy: formState.customProxy
    }

    await taskStore.createTask(args)
    message.success('任务创建成功')
    // 根据开关决定是否保留表单
    if (!keepFormAfterCreate.value) {
      resetForm()
    }
  } catch (err) {
    message.error(err.response?.data?.error || '创建任务失败')
  } finally {
    creating.value = false
  }
}

async function handleDelete(id) {
  try {
    await taskStore.deleteTask(id)
    message.success('删除成功')
  } catch {
    message.error('删除失败')
  }
}

// 保存任务为方案
async function saveAsProfile(task) {
  try {
    await profileStore.saveTaskAsProfile(task.id)
    message.success('方案保存成功')
  } catch {
    message.error('方案保存失败')
  }
}

// 加载方案到表单
function handleProfileChange(profileId) {
  if (!profileId) {
    // 清除方案时只重置配置字段，保留 url 和 outputName
    formState.threadCount = 32
    formState.retryCount = 15
    formState.headers = ''
    formState.baseUrl = ''
    formState.delAfterDone = true
    formState.binaryMerge = false
    formState.autoSelect = false
    formState.skipSegmentsCheck = false
    formState.concurrentDownload = false
    formState.decryptionEngine = 'MP4DECRYPT'
    formState.customArgs = ''
    formState.customProxy = ''
    return
  }

  const profile = profileStore.profiles.find(p => p.id === profileId)
  if (profile) {
    formState.threadCount = profile.thread_count
    formState.retryCount = profile.retry_count
    formState.headers = profile.headers || ''
    formState.baseUrl = profile.base_url || ''
    formState.delAfterDone = profile.del_after_done
    formState.binaryMerge = profile.binary_merge
    formState.autoSelect = profile.auto_select
    formState.skipSegmentsCheck = profile.skip_segments_check || false
    formState.concurrentDownload = profile.concurrent_download || false
    formState.decryptionEngine = profile.decryption_engine
    formState.customArgs = profile.custom_args || ''
    formState.customProxy = profile.custom_proxy || ''
    message.success(`已加载方案: ${profile.name}`)
  }
}

// 从方案管理弹窗加载方案
function loadProfileToForm(profile) {
  selectedProfileId.value = profile.id
  handleProfileChange(profile.id)
  showProfileManager.value = false
}

// 删除方案
async function handleDeleteProfile(id) {
  try {
    await profileStore.deleteProfile(id)
    message.success('方案删除成功')
    if (selectedProfileId.value === id) {
      selectedProfileId.value = null
    }
  } catch {
    message.error('方案删除失败')
  }
}

function viewLog(task) {
  logTaskId.value = task.id
  logModalVisible.value = true
}

onMounted(() => {
  // 使用 store 统一管理的轮询（单例模式）
  taskStore.startPolling()
  // 加载下载方案列表
  profileStore.fetchProfiles()
})
</script>

<style scoped>
.tasks-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.task-list {
  flex: 1;
}
</style>
