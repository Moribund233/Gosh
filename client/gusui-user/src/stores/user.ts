import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getUserInfo, updateUserInfo, type UserInfo } from '@/services/auth'

export const useUserStore = defineStore('user', () => {
  const userInfo = ref<UserInfo | null>(null)

  async function fetchUserInfo() {
    const res = await getUserInfo()
    userInfo.value = res
    return res
  }

  function setUserInfo(info: UserInfo) {
    userInfo.value = info
  }

  async function updateInfo(data: Partial<UserInfo>) {
    const res = await updateUserInfo(data)
    userInfo.value = res
    return res
  }

  function clearUser() {
    userInfo.value = null
  }

  return { userInfo, fetchUserInfo, setUserInfo, updateInfo, clearUser }
})
