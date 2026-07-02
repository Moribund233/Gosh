import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login, register, type LoginParams, type RegisterParams } from '@/services/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(uni.getStorageSync('token') || '')
  const isLoggedIn = computed(() => !!token.value)

  async function loginAction(params: LoginParams) {
    const res = await login(params)
    token.value = res.token
    uni.setStorageSync('token', res.token)
    uni.setStorageSync('userInfo', res.user)
    return res
  }

  async function registerAction(params: RegisterParams) {
    const res = await register(params)
    token.value = res.token
    uni.setStorageSync('token', res.token)
    uni.setStorageSync('userInfo', res.user)
    return res
  }

  function logout() {
    token.value = ''
    uni.removeStorageSync('token')
    uni.removeStorageSync('userInfo')
    uni.reLaunch({ url: '/pages/login/login' })
  }

  return { token, isLoggedIn, loginAction, registerAction, logout }
})
