import { defineStore } from 'pinia'
import { ref } from 'vue'
import { get, post, put, del } from '../api'

export const useProfileStore = defineStore('profile', () => {
  const profiles = ref([])
  const loading = ref(false)

  // 获取方案列表
  async function fetchProfiles(domain = '') {
    loading.value = true
    try {
      const url = domain ? `/profiles?domain=${encodeURIComponent(domain)}` : '/profiles'
      profiles.value = await get(url)
    } finally {
      loading.value = false
    }
  }

  // 创建方案
  async function createProfile(data) {
    const profile = await post('/profiles', data)
    profiles.value.unshift(profile)
    return profile
  }

  // 更新方案
  async function updateProfile(id, data) {
    const profile = await put(`/profiles/${id}`, data)
    const index = profiles.value.findIndex(p => p.id === id)
    if (index !== -1) {
      profiles.value[index] = profile
    }
    return profile
  }

  // 删除方案
  async function deleteProfile(id) {
    await del(`/profiles/${id}`)
    profiles.value = profiles.value.filter(p => p.id !== id)
  }

  // 根据域名获取方案
  async function getProfileByDomain(domain) {
    return await get('/profiles/by-domain', { domain })
  }

  // 将任务保存为方案
  // 返回 { profile, action, requested_name }：
  //   action = created | renamed | exists（exists 表示内容重复，未新增）
  async function saveTaskAsProfile(taskId) {
    const res = await post(`/tasks/${taskId}/save-as-profile`)
    const profile = res?.profile
    if (profile) {
      const index = profiles.value.findIndex(p => p.id === profile.id)
      if (index !== -1) {
        profiles.value[index] = profile
      } else {
        profiles.value.unshift(profile)
      }
    }
    return res
  }

  return {
    profiles,
    loading,
    fetchProfiles,
    createProfile,
    updateProfile,
    deleteProfile,
    getProfileByDomain,
    saveTaskAsProfile
  }
})
