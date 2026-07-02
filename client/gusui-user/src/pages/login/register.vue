<template>
  <view class="register-page">
    <view class="brand-section">
      <text class="brand-title">注册账号</text>
      <text class="brand-desc">注册后享受更多会员权益</text>
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

      <view class="form-group code-group">
        <text class="form-label">验证码</text>
        <view class="code-row">
          <input
            class="form-input code-input"
            v-model="code"
            type="text"
            maxlength="6"
            placeholder="请输入验证码"
            placeholder-class="placeholder"
          />
          <button class="code-btn" :disabled="codeSending || codeCountdown > 0" @tap="handleSendCode">
            <text v-if="codeCountdown > 0">{{ codeCountdown }}s</text>
            <text v-else-if="codeSending">发送中...</text>
            <text v-else>获取验证码</text>
          </button>
        </view>
      </view>

      <view class="form-group">
        <text class="form-label">设置密码</text>
        <input
          class="form-input"
          v-model="password"
          type="password"
          placeholder="请设置6-20位密码"
          placeholder-class="placeholder"
        />
      </view>

      <view class="form-group">
        <text class="form-label">确认密码</text>
        <input
          class="form-input"
          v-model="confirmPassword"
          type="password"
          placeholder="请再次输入密码"
          placeholder-class="placeholder"
        />
      </view>

      <button class="btn-primary" :disabled="!canRegister || loading" @tap="handleRegister">
        <text v-if="loading">注册中...</text>
        <text v-else>注 册</text>
      </button>

      <view class="login-row" @tap="navigateToLogin">
        <text class="login-link">已有账号？去登录</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useUserStore } from '@/stores/user'
import { sendCode } from '@/services/auth'

const authStore = useAuthStore()
const userStore = useUserStore()

const phone = ref('')
const code = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const codeSending = ref(false)
const codeCountdown = ref(0)

let countdownTimer: ReturnType<typeof setInterval> | null = null

const canRegister = computed(() =>
  phone.value.length === 11 &&
  code.value.length >= 4 &&
  password.value.length >= 6 &&
  password.value === confirmPassword.value
)

async function handleSendCode() {
  if (phone.value.length !== 11) {
    uni.showToast({ title: '请输入正确的手机号', icon: 'none' })
    return
  }
  codeSending.value = true
  try {
    await sendCode(phone.value)
    uni.showToast({ title: '验证码已发送', icon: 'success' })
    codeCountdown.value = 60
    countdownTimer = setInterval(() => {
      codeCountdown.value--
      if (codeCountdown.value <= 0) {
        if (countdownTimer) clearInterval(countdownTimer)
      }
    }, 1000)
  } catch {
    uni.showToast({ title: '发送失败，请重试', icon: 'none' })
  } finally {
    codeSending.value = false
  }
}

async function handleRegister() {
  if (!canRegister.value) return
  loading.value = true
  try {
    await authStore.registerAction({ phone: phone.value, password: password.value, code: code.value })
    if (countdownTimer) clearInterval(countdownTimer)
    userStore.setUserInfo(await userStore.fetchUserInfo())
    uni.switchTab({ url: '/pages/index/index' })
  } catch {
  } finally {
    loading.value = false
  }
}

function navigateToLogin() {
  uni.navigateBack()
}
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  background: #f7f8f5;
  padding: 40rpx 40rpx 40rpx;
}

.brand-section {
  margin-bottom: 48rpx;
  padding-top: 20rpx;
}

.brand-title {
  font-size: 40rpx;
  font-weight: 700;
  color: #2c2e2c;
}

.brand-desc {
  font-size: 24rpx;
  color: #9a9c98;
  margin-top: 8rpx;
  display: block;
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

.code-row {
  display: flex;
  gap: 20rpx;
}

.code-input {
  flex: 1;
}

.code-btn {
  width: 200rpx;
  height: 88rpx;
  background: #e8f0e6;
  color: #5b8c5a;
  border-radius: 16rpx;
  font-size: 24rpx;
  font-weight: 500;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  flex-shrink: 0;
}

.code-btn:active {
  opacity: 0.8;
}

.code-btn[disabled] {
  opacity: 0.5;
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
  margin-top: 20rpx;
}

.btn-primary:active {
  opacity: 0.85;
}

.btn-primary[disabled] {
  opacity: 0.5;
}

.login-row {
  text-align: center;
  margin-top: 32rpx;
}

.login-link {
  font-size: 24rpx;
  color: #5b8c5a;
}
</style>
