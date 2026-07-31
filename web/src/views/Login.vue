<template>
  <div class="login-page">
    <div class="login-backdrop" aria-hidden="true" />
    <div class="login-main">
      <div class="login-brand">
        <div class="brand-tile">
          <CloudDownloadOutlined />
        </div>
        <div class="brand-name">N_m3u8DL-RE</div>
        <div class="brand-sub">Web UI · m3u8 下载任务管理</div>
      </div>

      <a-config-provider :theme="cardTheme">
        <div class="login-card">
          <div class="login-card-title">登录</div>
          <a-form
            :model="formState"
            :rules="rules"
            layout="vertical"
            @finish="handleSubmit"
          >
            <a-form-item label="用户名" name="username">
              <a-input
                v-model:value="formState.username"
                placeholder="请输入用户名"
                size="large"
              >
                <template #prefix>
                  <UserOutlined class="field-icon" />
                </template>
              </a-input>
            </a-form-item>
            <a-form-item label="密码" name="password">
              <a-input-password
                v-model:value="formState.password"
                placeholder="请输入密码"
                size="large"
                @press-enter="handleSubmit"
              >
                <template #prefix>
                  <LockOutlined class="field-icon" />
                </template>
              </a-input-password>
            </a-form-item>
            <a-form-item style="margin-bottom: 0">
              <a-button
                type="primary"
                html-type="submit"
                :loading="loading"
                block
                size="large"
                class="login-button"
              >
                登录
              </a-button>
            </a-form-item>
          </a-form>
        </div>
      </a-config-provider>

      <div class="login-hint">默认用户名 admin · 默认密码 admin123</div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { UserOutlined, LockOutlined, CloudDownloadOutlined } from '@ant-design/icons-vue'
import { useUserStore } from '../stores/user'

const userStore = useUserStore()

const loading = ref(false)
const formState = reactive({
  username: '',
  password: ''
})

// 登录卡片固定使用亮色主题，与深色背景形成对比
const cardTheme = {
  token: {
    colorPrimary: '#6366f1',
    colorInfo: '#6366f1',
    colorBgContainer: '#ffffff',
    colorBgElevated: '#ffffff',
    colorText: '#1c2333',
    colorTextSecondary: '#5f6b85',
    colorBorder: '#dfe3ec',
    borderRadius: 10,
    borderRadiusLG: 12
  }
}

const rules = {
  username: [{ required: true, message: '请输入用户名' }],
  password: [{ required: true, message: '请输入密码' }]
}

async function handleSubmit() {
  loading.value = true
  try {
    await userStore.login(formState.username, formState.password)
    message.success('登录成功')
    // 使用 window.location 强制跳转，而不是 router.push
    window.location.href = '/'
  } catch (err) {
    message.error(err.response?.data?.error || err.message || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  overflow: hidden;
  background: #0b1020;
}

/* 渐变光斑 + 网格背景 */
.login-backdrop {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(640px 420px at 18% 12%, rgba(99, 102, 241, 0.32), transparent 62%),
    radial-gradient(560px 420px at 84% 78%, rgba(139, 92, 246, 0.26), transparent 62%),
    radial-gradient(420px 320px at 70% 15%, rgba(56, 189, 248, 0.14), transparent 65%);
}

.login-backdrop::after {
  content: '';
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.035) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.035) 1px, transparent 1px);
  background-size: 44px 44px;
  mask-image: radial-gradient(ellipse 70% 60% at 50% 45%, #000 30%, transparent 75%);
  -webkit-mask-image: radial-gradient(ellipse 70% 60% at 50% 45%, #000 30%, transparent 75%);
}

.login-main {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 380px;
  padding: 24px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.login-brand {
  text-align: center;
  margin-bottom: 28px;
}

.brand-tile {
  width: 56px;
  height: 56px;
  margin: 0 auto 14px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  color: #fff;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  box-shadow: 0 8px 24px rgba(99, 102, 241, 0.45);
}

.brand-name {
  font-size: 22px;
  font-weight: 700;
  color: #fff;
  letter-spacing: 0.3px;
}

.brand-sub {
  margin-top: 4px;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.55);
}

.login-card {
  width: 100%;
  padding: 28px 26px;
  border-radius: 16px;
  background: #fff;
  box-shadow: 0 24px 64px rgba(2, 6, 23, 0.55);
}

.login-card-title {
  margin-bottom: 20px;
  font-size: 17px;
  font-weight: 700;
  color: #1c2333;
}

.field-icon {
  color: #98a2b8;
}

.login-button {
  font-weight: 600;
  background: linear-gradient(135deg, #6366f1 0%, #7c6cf8 100%);
  border: none;
  box-shadow: 0 6px 16px rgba(99, 102, 241, 0.35);
}

.login-button:hover {
  opacity: 0.92;
}

.login-hint {
  margin-top: 18px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
}
</style>
