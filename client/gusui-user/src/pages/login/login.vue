<template>
  <view class="login-page">
    <view class="brand-section">
      <view class="brand-icon">
        <tui-icon name="explore" size="48" color="#5b8c5a"></tui-icon>
      </view>
      <text class="brand-name">谷穗</text>
      <text class="brand-sub">原产地干杂货 · 新鲜直达</text>
    </view>

    <view class="form-section">
      <view class="form-group">
        <text class="form-label">手机号</text>
        <input
          class="form-input"
          v-model="phone"
          type="text"
          maxlength="11"
          placeholder="请输入手机号"
          placeholder-class="placeholder"
        />
      </view>

      <view class="form-group">
        <text class="form-label">密码</text>
        <input
          class="form-input"
          v-model="password"
          type="password"
          placeholder="请输入密码"
          placeholder-class="placeholder"
        />
      </view>

      <view class="forgot-row" @tap="navigateToRegister">
        <text class="forgot-link">还没有账号？去注册</text>
      </view>

      <button class="btn-primary" :disabled="!phone || !password || loading" @tap="handleLogin">
        <text v-if="loading">登录中...</text>
        <text v-else>登 录</text>
      </button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useUserStore } from '@/stores/user'

const authStore = useAuthStore()
const userStore = useUserStore()

const phone = ref('')
const password = ref('')
const loading = ref(false)

async function handleLogin() {
  if (!phone.value || !password.value) return
  loading.value = true
  try {
    await authStore.loginAction({ phone: phone.value, password: password.value })
    userStore.setUserInfo(await userStore.fetchUserInfo())
    uni.switchTab({ url: '/pages/index/index' })
  } catch {
  } finally {
    loading.value = false
  }
}

function navigateToRegister() {
  uni.navigateTo({ url: '/pages/login/register' })
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: #f7f8f5;
  padding: 60rpx 40rpx 40rpx;
}

.brand-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 60rpx;
}

.brand-icon {
  width: 96rpx;
  height: 96rpx;
  background: #e8f0e6;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20rpx;
}

.brand-name {
  font-size: 44rpx;
  font-weight: 700;
  color: #2c2e2c;
  letter-spacing: 4rpx;
}

.brand-sub {
  font-size: 24rpx;
  color: #9a9c98;
  margin-top: 8rpx;
}

.form-section {
  padding: 0 8rpx;
}

.form-group {
  margin-bottom: 32rpx;
}

.form-label {
  font-size: 26rpx;
  color: #555755;
  font-weight: 500;
  margin-bottom: 12rpx;
  display: block;
}

.form-input {
  width: 100%;
  height: 88rpx;
  padding: 0 24rpx;
  background: #f0f2ed;
  border-radius: 16rpx;
  font-size: 28rpx;
  color: #2c2e2c;
  border: 2rpx solid transparent;
  box-sizing: border-box;
}

.form-input:focus {
  border-color: #5b8c5a;
  background: #ffffff;
}

.forgot-row {
  text-align: right;
  margin-bottom: 40rpx;
}

.forgot-link {
  font-size: 24rpx;
  color: #5b8c5a;
}

.btn-primary {
  width: 100%;
  height: 88rpx;
  background: #5b8c5a;
  color: #ffffff;
  border-radius: 16rpx;
  font-size: 30rpx;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  opacity: 1;
}

.btn-primary:active {
  opacity: 0.85;
}

.btn-primary[disabled] {
  opacity: 0.5;
}
</style>
