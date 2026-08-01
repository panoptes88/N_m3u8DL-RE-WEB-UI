<template>
  <div class="files-page">
    <PageHeader title="文件管理" subtitle="浏览、播放与下载已完成的文件">
      <a-input
        v-model:value="searchKeyword"
        placeholder="搜索文件名"
        allow-clear
        style="width: 200px"
      />
      <a-button @click="fetchFiles">
        <template #icon><ReloadOutlined /></template>
        刷新
      </a-button>
    </PageHeader>

    <a-card class="app-card files-card">
      <transition name="fade">
        <div class="batch-bar" v-if="selectedRowKeys.length > 0">
          <span class="batch-info">已选择 {{ selectedRowKeys.length }} 个文件</span>
          <a-button type="primary" danger size="small" @click="batchDelete">
            批量删除
          </a-button>
        </div>
      </transition>

      <a-table
        :columns="columns"
        :data-source="filteredFiles"
        :pagination="paginationConfig"
        :loading="loading"
        :row-selection="rowSelection"
        :row-key="record => record.name"
        :scroll="{ x: 600 }"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <a-space :size="8">
              <PlayCircleOutlined v-if="isVideoFile(record.name)" class="play-icon" @click="playVideo(record)" />
              <FileOutlined v-else class="file-icon" />
              <a-tooltip :title="record.name">
                <span class="file-name copyable-text" @click="copyToClipboard(record.name)">{{ record.name }}</span>
              </a-tooltip>
            </a-space>
          </template>
          <template v-if="column.key === 'size'">
            {{ formatSize(record.size) }}
          </template>
          <template v-if="column.key === 'modTime'">
            {{ record.modTime }}
          </template>
          <template v-if="column.key === 'action'">
            <a-space :size="4">
              <a-button size="small" type="text" @click="downloadFile(record.name)">
                下载
              </a-button>
              <a-popconfirm
                title="确定删除此文件？"
                @confirm="deleteFile(record.name)"
              >
                <a-button size="small" type="text" danger>删除</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 视频播放弹窗 -->
    <a-modal
      v-model:open="videoModalVisible"
      :title="currentVideoName"
      :footer="null"
      @cancel="closeVideo"
      class="video-modal"
      :width="800"
    >
      <div class="video-wrapper">
        <div id="xgplayer-container" class="xgplayer-container"></div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed, nextTick } from 'vue'
import { message } from 'ant-design-vue'
import { ReloadOutlined, PlayCircleOutlined, FileOutlined } from '@ant-design/icons-vue'
import { get, del } from '../api'
import Player from 'xgplayer'
import PageHeader from '../components/PageHeader.vue'

const loading = ref(false)
const files = ref([])
const searchKeyword = ref('')
const videoModalVisible = ref(false)
const currentVideoName = ref('')
const currentVideoUrl = ref('')
const selectedRowKeys = ref([])
const pageSize = ref(10)

// 按文件名搜索过滤
const filteredFiles = computed(() => {
  const kw = searchKeyword.value.trim().toLowerCase()
  if (!kw) return files.value
  return files.value.filter(f => f.name.toLowerCase().includes(kw))
})

// 西瓜播放器实例
let xgPlayer = null

const videoExtensions = ['.mp4', '.mkv', '.avi', '.mov', '.webm', '.flv', '.wmv', '.m4v', '.3gp']

const paginationConfig = computed(() => ({
  pageSize: pageSize.value,
  showSizeChanger: true,
  pageSizeOptions: ['10', '20', '50', '100', '200'],
  showTotal: (total) => `共 ${total} 个文件`
}))

const columns = [
  {
    title: '文件名',
    dataIndex: 'name',
    key: 'name',
    ellipsis: true,
    sorter: (a, b) => a.name.localeCompare(b.name),
    width: 150
  },
  {
    title: '大小',
    dataIndex: 'size',
    key: 'size',
    width: 90,
    sorter: (a, b) => a.size - b.size
  },
  {
    title: '修改时间',
    dataIndex: 'modTime',
    key: 'modTime',
    width: 140,
    sorter: (a, b) => new Date(a.modTime) - new Date(b.modTime),
    responsive: ['lg']
  },
  { title: '操作', key: 'action', width: 110 }
]

const rowSelection = computed(() => ({
  selectedRowKeys: selectedRowKeys.value,
  onChange: (keys) => {
    selectedRowKeys.value = keys
  }
}))

function formatSize(bytes) {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function handleTableChange(paginationInfo) {
  pageSize.value = paginationInfo.pageSize
}

function isVideoFile(filename) {
  const ext = filename.substring(filename.lastIndexOf('.')).toLowerCase()
  return videoExtensions.includes(ext)
}

async function fetchFiles() {
  loading.value = true
  try {
    files.value = await get('/files')
    selectedRowKeys.value = []
  } catch (err) {
    message.error('获取文件列表失败')
  } finally {
    loading.value = false
  }
}

function downloadFile(name) {
  window.open(`/api/files/download?name=${encodeURIComponent(name)}`, '_blank')
}

// 点击复制到粘贴板
async function copyToClipboard(text) {
  try {
    if (navigator.clipboard) {
      await navigator.clipboard.writeText(text)
      message.success('已复制到粘贴板')
      return
    }
  } catch {
    // 忽略，回退
  }
  try {
    const textarea = document.createElement('textarea')
    textarea.value = text
    textarea.style.position = 'fixed'
    textarea.style.opacity = '0'
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
    message.success('已复制到粘贴板')
  } catch {
    message.error('复制失败')
  }
}

async function playVideo(record) {
  currentVideoName.value = record.name
  currentVideoUrl.value = `/api/files/download?name=${encodeURIComponent(record.name)}`
  videoModalVisible.value = true

  // 修改页面标题为文件名（用于 QQ 浏览器嗅探）
  document.title = record.name

  // 等待 DOM 更新后初始化播放器
  await nextTick()
  initXgPlayer(currentVideoUrl.value, record.name)
  // 尝试加载同名字幕
  loadSubtitle(record.name)
}

// SRT 字幕转 VTT
function srtToVtt(srt) {
  return 'WEBVTT\n\n' + srt
    .replace(/\r\n/g, '\n')
    .replace(/\r/g, '\n')
    .replace(/(\d{2}:\d{2}:\d{2}),(\d{3})/g, '$1.$2')
}

// 加载同名字幕（.srt/.vtt），用 HTML5 track 原生渲染
async function loadSubtitle(videoName) {
  const baseName = videoName.replace(/\.[^.]+$/, '')
  // 匹配同名字幕：video.mp4 -> video.srt / video.en.srt / video.zh.srt 等
  const subFiles = files.value.filter(f =>
    f.name.startsWith(baseName + '.') && /\.(srt|vtt)$/i.test(f.name)
  )
  // 优先中/英文字幕，否则取第一个
  const subFile = subFiles.find(f => /\.(zh|en)\./i.test(f.name)) || subFiles[0]
  if (!subFile) return

  try {
    const subUrl = `/api/files/download?name=${encodeURIComponent(subFile.name)}`
    const res = await fetch(subUrl)
    let text = await res.text()
    if (/\.srt$/i.test(subFile.name)) {
      text = srtToVtt(text)
    }
    // 用 HTML5 <track> 加载字幕，浏览器原生渲染
    const blob = new Blob([text], { type: 'text/vtt' })
    const vttUrl = URL.createObjectURL(blob)
    const video = document.querySelector('#xgplayer-container video')
    if (!video) {
      console.warn('[字幕] 未找到 video 元素')
      return
    }
    // 清除旧 track
    video.querySelectorAll('track').forEach(t => t.remove())
    const track = document.createElement('track')
    track.kind = 'subtitles'
    track.label = '字幕'
    track.srclang = 'zh'
    track.src = vttUrl
    track.default = true
    video.appendChild(track)
    // 字幕加载后开启显示
    track.addEventListener('load', () => {
      if (video.textTracks && video.textTracks.length > 0) {
        video.textTracks[0].mode = 'showing'
      }
    })
  } catch (err) {
    console.error('字幕加载失败:', err)
    message.warning('字幕加载失败，请查看控制台')
  }
}

function initXgPlayer(url, title = '') {
  // 如果已存在播放器，先销毁
  if (xgPlayer) {
    xgPlayer.destroy()
    xgPlayer = null
  }

  const savedVolume = parseFloat(localStorage.getItem('xg_volume') || '1')
  xgPlayer = new Player({
    id: 'xgplayer-container',
    url: url,
    playbackRate: [0.5, 0.75, 1, 1.25, 1.5, 2],
    defaultPlaybackRate: 1,
    fluid: true,
    volume: savedVolume,
    maxVolume: 1,
    ignores: ['quality'],
    closeVideoClick: false,
    title: title,
    error: () => {
      message.error('视频加载失败，请检查文件是否完整')
    }
  })

  // 音量记忆：变化时保存
  xgPlayer.on('volumechange', () => {
    localStorage.setItem('xg_volume', String(xgPlayer.volume))
  })

  // 播放进度记忆：恢复上次播放位置
  const progressKey = `xg_progress_${title}`
  const savedTime = parseFloat(localStorage.getItem(progressKey) || '0')
  if (savedTime > 0) {
    xgPlayer.once('canplay', () => {
      xgPlayer.currentTime = savedTime
    })
  }
  xgPlayer.on('timeupdate', () => {
    localStorage.setItem(progressKey, String(xgPlayer.currentTime))
  })
  xgPlayer.on('ended', () => {
    localStorage.removeItem(progressKey)
  })
}

function closeVideo() {
  if (xgPlayer) {
    // 关闭前保存最后进度
    if (currentVideoName.value && xgPlayer.currentTime > 0) {
      localStorage.setItem(`xg_progress_${currentVideoName.value}`, String(xgPlayer.currentTime))
    }
    xgPlayer.destroy()
    xgPlayer = null
  }
  // 恢复页面标题
  document.title = 'N_m3u8DL-RE Web UI'
  currentVideoName.value = ''
  currentVideoUrl.value = ''
}

// 组件销毁时清理播放器
onUnmounted(() => {
  if (xgPlayer) {
    xgPlayer.destroy()
    xgPlayer = null
  }
})

async function deleteFile(name) {
  try {
    await del(`/files/${encodeURIComponent(name)}`)
    message.success('删除成功')
    fetchFiles()
  } catch (err) {
    message.error(err.response?.data?.error || '删除失败')
  }
}

async function batchDelete() {
  if (selectedRowKeys.value.length === 0) {
    return
  }

  const names = selectedRowKeys.value.filter(name => {
    const file = files.value.find(f => f.name === name)
    return file && !file.isDir
  })

  if (names.length === 0) {
    message.warning('没有可删除的文件')
    return
  }

  try {
    let success = 0
    let failed = 0
    for (const name of names) {
      try {
        await del(`/files/${encodeURIComponent(name)}`)
        success++
      } catch {
        failed++
      }
    }

    if (success > 0) {
      message.success(`成功删除 ${success} 个文件`)
    }
    if (failed > 0) {
      message.warning(`删除失败 ${failed} 个文件`)
    }

    selectedRowKeys.value = []
    fetchFiles()
  } catch (err) {
    message.error('批量删除失败')
  }
}

onMounted(() => {
  fetchFiles()
})
</script>

<style scoped>
.files-page {
  display: flex;
  flex-direction: column;
}

.files-card :deep(.ant-card-body) {
  padding-top: 12px;
}

.batch-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 14px;
  margin-bottom: 12px;
  border-radius: 10px;
  background: rgba(99, 102, 241, 0.08);
  border: 1px solid rgba(99, 102, 241, 0.2);
}

.batch-info {
  font-size: 13px;
  color: var(--text-1);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.play-icon {
  color: var(--brand);
  cursor: pointer;
  font-size: 17px;
  transition: transform 0.2s;
}

.play-icon:hover {
  transform: scale(1.15);
}

.file-icon {
  color: var(--text-3);
  font-size: 16px;
}

.file-name {
  color: var(--text-1);
}

.copyable-text {
  cursor: pointer;
  display: inline-block;
  max-width: 100%;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  vertical-align: bottom;
}

.copyable-text:hover {
  color: var(--brand);
}

.video-wrapper {
  position: relative;
}

.xgplayer-container {
  width: 100%;
  max-height: 450px;
}

.video-modal :deep(.ant-modal-body) {
  padding: 8px;
}

@media (max-width: 576px) {
  .files-page :deep(.ant-table-cell) {
    padding: 8px !important;
  }
}
</style>
