<template>
  <view class="profile-page">
    <view class="nav-bar">
      <text class="nav-title">个人中心</text>
      <view class="nav-right" @tap="goSettings">
        <tui-icon name="setup" size="36" color="#555755"></tui-icon>
      </view>
    </view>

    <scroll-view class="content" scroll-y>
      <view class="profile-header" @tap="handleLogin">
        <view class="avatar" v-if="userStore.userInfo">
          <text class="avatar-text">{{ userStore.userInfo.nickname?.[0] || 'U' }}</text>
        </view>
        <view class="avatar" v-else>
          <tui-icon name="people" size="40" color="#9a9c98"></tui-icon>
        </view>
        <view class="header-info">
          <text class="header-name">{{ userStore.userInfo?.nickname || '未登录' }}</text>
          <text class="header-detail" v-if="userStore.userInfo">
            会员ID: {{ userStore.userInfo.id }} · 欢迎回来
          </text>
          <text class="header-detail" v-else>
            登录后享受更多权益
          </text>
        </view>
        <tui-icon name="arrowright" size="28" color="#bcbeb9"></tui-icon>
      </view>

      <view class="order-section" @tap="goOrders">
        <view class="order-section-header">
          <text class="order-section-title">我的订单</text>
          <text class="order-section-more">查看全部 &gt;</text>
        </view>
        <view class="order-stats">
          <view class="stat-item" @tap.stop="goOrders('unpaid')">
            <tui-icon name="time" size="40" color="#555755"></tui-icon>
            <text class="stat-label">待付款</text>
          </view>
          <view class="stat-item" @tap.stop="goOrders('undelivered')">
            <tui-icon name="transport" size="40" color="#555755"></tui-icon>
            <text class="stat-label">待发货</text>
          </view>
          <view class="stat-item" @tap.stop="goOrders('delivering')">
            <tui-icon name="location" size="40" color="#555755"></tui-icon>
            <text class="stat-label">待收货</text>
          </view>
          <view class="stat-item" @tap.stop="goOrders('delivered')">
            <tui-icon name="star" size="40" color="#555755"></tui-icon>
            <text class="stat-label">待评价</text>
          </view>
          <view class="stat-item" @tap.stop="goOrders('refund')">
            <tui-icon name="exchange" size="40" color="#555755"></tui-icon>
            <text class="stat-label">退换</text>
          </view>
        </view>
      </view>

      <view class="menu-group">
        <view class="menu-item" @tap="goAddress">
          <tui-icon name="location" size="36" color="#555755"></tui-icon>
          <text class="menu-label">收货地址</text>
          <tui-icon name="arrowright" size="28" color="#bcbeb9"></tui-icon>
        </view>
        <view class="menu-item" @tap="goCoupon">
          <tui-icon name="coupon" size="36" color="#555755"></tui-icon>
          <text class="menu-label">我的卡券</text>
          <text class="menu-badge" v-if="couponCount">3张</text>
          <tui-icon name="arrowright" size="28" color="#bcbeb9"></tui-icon>
        </view>
        <view class="menu-item" @tap="goFavorites">
          <tui-icon name="like" size="36" color="#555755"></tui-icon>
          <text class="menu-label">我的收藏</text>
          <tui-icon name="arrowright" size="28" color="#bcbeb9"></tui-icon>
        </view>
        <view class="menu-item" @tap="goHistory">
          <tui-icon name="time" size="36" color="#555755"></tui-icon>
          <text class="menu-label">浏览记录</text>
          <tui-icon name="arrowright" size="28" color="#bcbeb9"></tui-icon>
        </view>
      </view>

      <view class="menu-group">
        <view class="menu-item" @tap="goPoints">
          <tui-icon name="wealth" size="36" color="#555755"></tui-icon>
          <text class="menu-label">积分商城</text>
          <text class="menu-badge" v-if="pointsBalance">{{ pointsBalance }}分</text>
          <tui-icon name="arrowright" size="28" color="#bcbeb9"></tui-icon>
        </view>
        <view class="menu-item" @tap="goService">
          <tui-icon name="kefu" size="36" color="#555755"></tui-icon>
          <text class="menu-label">联系客服</text>
          <tui-icon name="arrowright" size="28" color="#bcbeb9"></tui-icon>
        </view>
        <view class="menu-item" @tap="goSettings">
          <tui-icon name="setup" size="36" color="#555755"></tui-icon>
          <text class="menu-label">设置</text>
          <tui-icon name="arrowright" size="28" color="#bcbeb9"></tui-icon>
        </view>
      </view>
    </scroll-view>
    <AppTabbar :current="3" />
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import AppTabbar from '@/components/AppTabbar.vue'
import { useAuthStore } from '@/stores/auth'
import { useUserStore } from '@/stores/user'

const authStore = useAuthStore()
const userStore = useUserStore()
const couponCount = ref(0)
const pointsBalance = ref(0)

onMounted(() => {
  if (authStore.isLoggedIn) {
    userStore.fetchUserInfo()
  }
})

function handleLogin() {
  if (!authStore.isLoggedIn) {
    uni.navigateTo({ url: '/pages/login/login' })
  }
}

function goOrders(status?: string) {
  const url = status
    ? `/pages/sub-package-order/orders/orders?status=${status}`
    : '/pages/sub-package-order/orders/orders'
  uni.navigateTo({ url })
}

function goAddress() {
  if (!authStore.isLoggedIn) { uni.navigateTo({ url: '/pages/login/login' }); return }
  uni.navigateTo({ url: '/pages/address/address' })
}

function goCoupon() {
  if (!authStore.isLoggedIn) { uni.navigateTo({ url: '/pages/login/login' }); return }
  uni.navigateTo({ url: '/pages/coupon/coupon' })
}

function goPoints() {
  if (!authStore.isLoggedIn) { uni.navigateTo({ url: '/pages/login/login' }); return }
  uni.navigateTo({ url: '/pages/points/points' })
}

function goFavorites() {
  if (!authStore.isLoggedIn) { uni.navigateTo({ url: '/pages/login/login' }); return }
  uni.showToast({ title: '功能开发中', icon: 'none' })
}

function goHistory() {
  if (!authStore.isLoggedIn) { uni.navigateTo({ url: '/pages/login/login' }); return }
  uni.showToast({ title: '功能开发中', icon: 'none' })
}

function goService() {
  uni.showToast({ title: '请联系客服电话 400-888-8888', icon: 'none' })
}

function goSettings() {
  if (!authStore.isLoggedIn) { uni.navigateTo({ url: '/pages/login/login' }); return }
  uni.showActionSheet({
    itemList: ['退出登录'],
    success: (res) => {
      if (res.tapIndex === 0) {
        authStore.logout()
        userStore.clearUser()
      }
    },
  })
}
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  background: #f7f8f5;
}

.nav-bar {
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(180deg, #e8f0e6, #f7f8f5);
  padding: 0 32rpx;
  position: relative;
}

.nav-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #2c2e2c;
}

.nav-right {
  position: absolute;
  right: 32rpx;
  top: 50%;
  transform: translateY(-50%);
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.content { height: calc(100vh - 88rpx); }

.profile-header {
  background: linear-gradient(180deg, #e8f0e6 0%, #ffffff 100%);
  padding: 40rpx 32rpx;
  display: flex;
  align-items: center;
  gap: 24rpx;
}

.avatar {
  width: 112rpx;
  height: 112rpx;
  border-radius: 50%;
  background: #5b8c5a;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.avatar-text {
  font-size: 44rpx;
  font-weight: 700;
  color: #ffffff;
}

.header-info { flex: 1; }

.header-name {
  font-size: 36rpx;
  font-weight: 700;
  color: #2c2e2c;
  display: block;
  margin-bottom: 4rpx;
}

.header-detail {
  font-size: 24rpx;
  color: #9a9c98;
}

.order-section {
  background: #ffffff;
  margin: 0 0 16rpx;
  border-radius: 0;
}

.order-section-header {
  display: flex;
  align-items: center;
  padding: 28rpx 32rpx 16rpx;
  border-bottom: 2rpx solid #eef0eb;
}

.order-section-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #2c2e2c;
  flex: 1;
}

.order-section-more {
  font-size: 24rpx;
  color: #9a9c98;
}

.order-stats {
  display: flex;
  padding: 20rpx 0;
}

.stat-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8rpx;
}

.stat-label {
  font-size: 22rpx;
  color: #9a9c98;
}

.menu-group {
  background: #ffffff;
  margin-bottom: 16rpx;
  border-radius: 16rpx;
  margin-left: 32rpx;
  margin-right: 32rpx;
  overflow: hidden;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 28rpx 24rpx;
  border-bottom: 2rpx solid #eef0eb;
  gap: 20rpx;
}

.menu-item:last-child {
  border-bottom: none;
}

.menu-label {
  flex: 1;
  font-size: 28rpx;
  color: #2c2e2c;
}

.menu-badge {
  font-size: 22rpx;
  color: #5b8c5a;
  background: #e8f0e6;
  padding: 4rpx 16rpx;
  border-radius: 999rpx;
}
</style>
