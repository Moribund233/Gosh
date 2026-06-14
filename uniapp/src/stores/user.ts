import { defineStore } from 'pinia'

interface UserInfo {
  id: number
  nickname: string
  avatar: string
  phone: string
}

export const useUserStore = defineStore('user', {
  state: () => ({
    token: '' as string,
    userInfo: null as UserInfo | null,
  }),
  getters: {
    isLoggedIn: (state) => !!state.token,
  },
  actions: {
    setToken(token: string) {
      this.token = token
      uni.setStorageSync('token', token)
    },
    setUserInfo(info: UserInfo) {
      this.userInfo = info
    },
    logout() {
      this.token = ''
      this.userInfo = null
      uni.removeStorageSync('token')
      uni.removeStorageSync('userInfo')
    },
  },
})
