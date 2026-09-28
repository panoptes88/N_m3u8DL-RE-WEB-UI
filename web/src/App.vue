<template>
  <a-config-provider :theme="themeConfig" :locale="zhCN">
    <div class="app-container">
      <template v-if="userStore.isLoggedIn">
        <a-layout class="layout">
          <!-- PC端左侧导航栏 (移动端隐藏)；theme="light" 使背景走 colorBgContainer，亮暗自适应 -->
          <a-layout-sider
            v-model:collapsed="appStore.collapsed"
            :trigger="null"
            collapsible
            theme="light"
            class="sider pc-sider"
            :class="{ 'sider-collapsed-hidden': appStore.collapsed && isMobile }"
            :width="208"
            :collapsed-width="72"
          >
            <BrandLogo :collapsed="appStore.collapsed" class="sider-logo" />
            <a-menu
              :selectedKeys="selectedKeys"
              mode="inline"
              :items="menuItems"
              @click="handleMenuClick"
            />
          </a-layout-sider>

          <a-layout>
            <!-- 顶部栏 -->
            <a-layout-header class="header">
              <div class="header-left">
                <!-- 移动端菜单按钮 -->
                <a-button
                  type="text"
                  class="trigger mobile-trigger"
                  aria-label="打开导航菜单"
                  @click="mobileMenuVisible = true"
                >
                  <template #icon><menu-unfold-outlined /></template>
                </a-button>
                <!-- PC端收缩按钮 -->
                <a-button
                  type="text"
                  class="trigger pc-trigger"
                  :aria-label="appStore.collapsed ? '展开侧边栏' : '收起侧边栏'"
                  @click="appStore.toggleCollapsed"
                >
                  <template #icon>
                    <menu-unfold-outlined v-if="appStore.collapsed" />
                    <menu-fold-outlined v-else />
                  </template>
                </a-button>
              </div>

              <div class="header-right">
                <a-space :size="8">
                  <!-- 主题切换：亮色 / 暗色 / 跟随系统 -->
                  <a-dropdown placement="bottomRight">
                    <a-button
                      type="text"
                      class="header-action"
                      :aria-label="`主题：${themeLabel}，点击切换`"
                      :title="`主题：${themeLabel}`"
                    >
                      <template #icon>
                        <bulb-filled v-if="appStore.resolvedTheme === 'dark'" />
                        <bulb-outlined v-else />
                      </template>
                    </a-button>
                    <template #overlay>
                      <a-menu :selectedKeys="[appStore.theme]" @click="handleThemeMenuClick">
                        <a-menu-item key="light">
                          <bulb-outlined />
                          <span style="margin-left: 8px;">亮色</span>
                        </a-menu-item>
                        <a-menu-item key="dark">
                          <bulb-filled />
                          <span style="margin-left: 8px;">暗色</span>
                        </a-menu-item>
                        <a-menu-item key="system">
                          <desktop-outlined />
                          <span style="margin-left: 8px;">跟随系统</span>
                        </a-menu-item>
                      </a-menu>
                    </template>
                  </a-dropdown>

                  <!-- 用户下拉菜单 -->
                  <a-dropdown>
                    <button type="button" class="user-dropdown" aria-label="用户菜单">
                      <a-avatar :size="30" class="user-avatar">
                        <template #icon><user-outlined /></template>
                      </a-avatar>
                      <span class="username">{{ userStore.username }}</span>
                    </button>
                    <template #overlay>
                      <a-menu @click="handleUserMenuClick">
                        <a-menu-item key="change-password">
                          <key-outlined />
                          <span style="margin-left: 8px;">修改密码</span>
                        </a-menu-item>
                        <a-menu-divider />
                        <a-menu-item key="logout">
                          <logout-outlined />
                          <span style="margin-left: 8px;">退出登录</span>
                        </a-menu-item>
                      </a-menu>
                    </template>
                  </a-dropdown>
                </a-space>
              </div>
            </a-layout-header>

            <!-- 内容区 -->
            <a-layout-content class="content">
              <router-view />
            </a-layout-content>
          </a-layout>
        </a-layout>

        <!-- 移动端抽屉式侧边栏 -->
        <a-drawer
          v-model:open="mobileMenuVisible"
          placement="left"
          :width="208"
          :closable="false"
          :show-header="false"
          class="mobile-drawer"
          :body-style="{ padding: 0 }"
        >
          <BrandLogo class="sider-logo" />
          <a-menu
            :selectedKeys="selectedKeys"
            mode="inline"
            :items="menuItems"
            @click="handleMobileMenuClick"
          />
        </a-drawer>
      </template>
      <template v-else>
        <router-view />
      </template>

      <!-- 修改密码弹窗 -->
      <a-modal
        v-model:open="changePasswordVisible"
        title="修改密码"
        :confirm-loading="changePasswordLoading"
        @ok="handleChangePasswordOk"
        @cancel="handleChangePasswordCancel"
      >
        <a-alert
          type="warning"
          show-icon
          message="修改成功后需要重新登录"
          description="本系统为单用户模式且没有找回密码入口，请务必确认两次输入一致。"
          style="margin-bottom: 16px"
        />
        <a-form
          ref="changePasswordFormRef"
          :model="changePasswordForm"
          :rules="changePasswordRules"
          layout="vertical"
        >
          <a-form-item label="新密码" name="newPassword">
            <a-input-password
              v-model:value="changePasswordForm.newPassword"
              placeholder="请输入新密码"
            />
          </a-form-item>
          <a-form-item label="确认新密码" name="confirmPassword">
            <a-input-password
              v-model:value="changePasswordForm.confirmPassword"
              placeholder="请再次输入新密码"
              @pressEnter="handleChangePasswordOk"
            />
          </a-form-item>
        </a-form>
      </a-modal>
    </div>
  </a-config-provider>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted, h } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import zhCN from 'ant-design-vue/es/locale/zh_CN'
import { useUserStore } from './stores/user'
import { useAppStore } from './stores/app'
import { getThemeConfig } from './theme'
import { post } from './api'
import BrandLogo from './components/BrandLogo.vue'
import {
  DashboardOutlined,
  CloudDownloadOutlined,
  FolderOutlined,
  UserOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  LogoutOutlined,
  KeyOutlined,
  BulbOutlined,
  BulbFilled,
  DesktopOutlined
} from '@ant-design/icons-vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const appStore = useAppStore()

const mobileMenuVisible = ref(false)
const isMobile = ref(window.innerWidth <= 768)

const themeConfig = computed(() => getThemeConfig(appStore.resolvedTheme))

// 主题模式的中文名，用于按钮提示与无障碍标签
const THEME_LABELS = { light: '亮色', dark: '暗色', system: '跟随系统' }
const themeLabel = computed(() => THEME_LABELS[appStore.theme] || '亮色')

function handleThemeMenuClick({ key }) {
  appStore.setTheme(key)
}

const handleResize = () => {
  isMobile.value = window.innerWidth <= 768
}

// 修改密码相关
const changePasswordVisible = ref(false)
const changePasswordLoading = ref(false)
const changePasswordFormRef = ref(null)
const changePasswordForm = reactive({
  newPassword: '',
  confirmPassword: ''
})

// 单用户系统没有找回入口，输错即锁死，因此必须二次确认
const changePasswordRules = {
  newPassword: [{ required: true, message: '请输入新密码' }],
  confirmPassword: [
    { required: true, message: '请再次输入新密码' },
    {
      validator: (_rule, value) =>
        value === changePasswordForm.newPassword
          ? Promise.resolve()
          : Promise.reject('两次输入的密码不一致')
    }
  ]
}

const menuItems = [
  {
    key: 'dashboard',
    icon: () => h(DashboardOutlined),
    label: '首页'
  },
  {
    key: 'tasks',
    icon: () => h(CloudDownloadOutlined),
    label: '下载任务'
  },
  {
    key: 'files',
    icon: () => h(FolderOutlined),
    label: '文件管理'
  }
]

// 菜单选中态跟随路由（刷新、前进后退、直输 URL 时保持同步）
// 用 computed 单向跟随 route.name，避免硬编码初始值导致刷新时高亮先落在"首页"再跳转
const selectedKeys = computed(() => {
  const name = route.name
  if (typeof name === 'string') {
    const key = name.toLowerCase()
    if (menuItems.some(item => item.key === key)) {
      return [key]
    }
  }
  return []
})

function handleMenuClick({ key }) {
  router.push({ name: key.charAt(0).toUpperCase() + key.slice(1) })
}

function handleMobileMenuClick({ key }) {
  mobileMenuVisible.value = false
  router.push({ name: key.charAt(0).toUpperCase() + key.slice(1) })
}

function handleUserMenuClick({ key }) {
  if (key === 'change-password') {
    changePasswordVisible.value = true
    changePasswordForm.newPassword = ''
    changePasswordForm.confirmPassword = ''
  } else if (key === 'logout') {
    userStore.logout().then(() => {
      router.push('/login')
    })
  }
}

async function handleChangePasswordOk() {
  try {
    await changePasswordFormRef.value.validate()
  } catch {
    return
  }

  changePasswordLoading.value = true
  try {
    await post('/auth/change-password', {
      new_password: changePasswordForm.newPassword
    })
    message.success('密码修改成功')
    changePasswordVisible.value = false
    userStore.logout().then(() => {
      router.push('/login')
    })
  } catch (err) {
    message.error(err.response?.data?.error || '修改密码失败')
  } finally {
    changePasswordLoading.value = false
  }
}

function handleChangePasswordCancel() {
  changePasswordVisible.value = false
  changePasswordForm.newPassword = ''
  changePasswordForm.confirmPassword = ''
  changePasswordFormRef.value?.clearValidate()
}

onMounted(async () => {
  window.addEventListener('resize', handleResize)
  const isLoggedIn = await userStore.checkLogin()
  if (!isLoggedIn && route.name !== 'Login') {
    router.push('/login')
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.app-container {
  min-height: 100vh;
  background: var(--bg-body);
}

.layout {
  min-height: 100vh;
}

/* PC端侧边栏 */
.sider {
  border-inline-end: 1px solid var(--border-soft);
}

.sider-logo {
  border-bottom: 1px solid var(--border-soft);
}

.sider :deep(.ant-menu) {
  border-inline-end: none;
  padding: 8px 0;
}

/* 顶部栏：吸顶 + 毛玻璃 */
.header {
  position: sticky;
  top: 0;
  z-index: 20;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--border-soft);
  backdrop-filter: saturate(180%) blur(10px);
  -webkit-backdrop-filter: saturate(180%) blur(10px);
}

.header-left {
  display: flex;
  align-items: center;
}

.trigger {
  font-size: 17px;
  color: var(--text-2);
}

.trigger:hover {
  color: var(--brand) !important;
  background: rgba(99, 102, 241, 0.08) !important;
}

.header-right {
  display: flex;
  align-items: center;
}

.header-action {
  color: var(--text-2);
}

/* 触发下拉的用户区：用真实 button 以获得键盘可达性，同时清掉默认按钮外观 */
.user-dropdown {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 8px;
  border: none;
  background: none;
  border-radius: var(--radius-sm);
  font: inherit;
  color: inherit;
  transition: background 0.2s;
}

.user-dropdown:hover {
  background: rgba(99, 102, 241, 0.08);
}

.user-avatar {
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
}

.username {
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 13px;
  color: var(--text-1);
}

.content {
  padding: 24px;
  margin: 0;
  overflow: auto;
}

/* PC端触发器 */
.pc-trigger {
  display: inline-flex;
}

/* 移动端触发器 */
.mobile-trigger {
  display: none;
}

/* 移动端响应式 */
@media (max-width: 768px) {
  .pc-sider {
    display: none;
  }

  .pc-trigger {
    display: none;
  }

  .mobile-trigger {
    display: inline-flex;
  }

  .header {
    padding: 0 12px;
  }

  .content {
    padding: 16px;
  }

  .username {
    display: none;
  }
}

/* 移动端侧边栏完全隐藏 - 使用独立类避免 !important */
.sider-collapsed-hidden {
  width: 0;
  min-width: 0;
  max-width: 0;
  flex: 0 0 0;
  overflow: hidden;
}
</style>
