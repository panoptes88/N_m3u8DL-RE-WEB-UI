import axios from 'axios'
import { message } from 'ant-design-vue'

const api = axios.create({
  baseURL: '/api',
  timeout: 30000,
  withCredentials: true
})

// 会话过期只处理一次，避免并发请求同时弹出多个提示、触发多次跳转
let sessionExpiredHandled = false

// 响应拦截器 - 返回响应数据；统一处理会话过期
api.interceptors.response.use(
  response => response.data,
  error => {
    const status = error.response?.status
    const url = error.config?.url || ''

    // 这两类 401 由调用方自行处理，不在此跳转：
    //   /auth/login -> 密码错误，登录页自己提示
    //   /user       -> 启动时的登录态探测，App.vue 已有跳转逻辑
    const handledByCaller = url.includes('/auth/login') || url.includes('/user')

    if (status === 401 && !handledByCaller) {
      localStorage.removeItem('username')

      const onLoginPage = window.location.pathname === '/login'
      if (!onLoginPage && !sessionExpiredHandled) {
        sessionExpiredHandled = true
        message.error('登录已过期，请重新登录')
        const redirect = encodeURIComponent(
          window.location.pathname + window.location.search
        )
        // 整页跳转以彻底重置内存中的登录态与轮询
        setTimeout(() => {
          window.location.href = `/login?redirect=${redirect}`
        }, 600)
      }
    }

    return Promise.reject(error)
  }
)

export function get(url, params) {
  return api.get(url, { params })
}

export function post(url, data) {
  return api.post(url, data)
}

export function del(url) {
  return api.delete(url)
}

export function put(url, data) {
  return api.put(url, data)
}

export default api
