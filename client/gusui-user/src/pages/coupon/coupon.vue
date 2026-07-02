<template>
  <view class="coupon-page">
    <scroll-view class="content" scroll-y>
      <view class="coupon-list">
        <view class="coupon-card" v-for="item in list" :key="item.id">
          <view class="coupon-left">
            <text class="coupon-value">{{ item.type === 'fixed' ? ('¥' + item.value) : (item.value + '%') }}</text>
            <text class="coupon-condition">满{{ item.min_amount }}元可用</text>
          </view>
          <view class="coupon-right">
            <text class="coupon-name">{{ item.name }}</text>
            <text class="coupon-date">{{ item.end_at }} 到期</text>
            <text class="coupon-status" :class="{ used: item.used }">{{ item.used ? '已使用' : '未使用' }}</text>
          </view>
        </view>
      </view>

      <view class="empty-state" v-if="list.length === 0 && !loading">
        <tui-icon name="coupon" size="64" color="#bcbeb9"></tui-icon>
        <text class="empty-text">暂无优惠券</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getCoupons, type Coupon } from '@/services/coupon'

const list = ref<Coupon[]>([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    list.value = await getCoupons()
  } catch {
    list.value = []
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.coupon-page {
  min-height: 100vh;
  background: #f7f8f5;
}

.content { height: 100vh; }

.coupon-list {
  padding: 24rpx 32rpx;
}

.coupon-card {
  display: flex;
  background: #ffffff;
  border-radius: 16rpx;
  overflow: hidden;
  margin-bottom: 20rpx;
}

.coupon-left {
  width: 200rpx;
  background: #e8f0e6;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 32rpx 16rpx;
}

.coupon-value {
  font-size: 44rpx;
  font-weight: 700;
  color: #5b8c5a;
}

.coupon-condition {
  font-size: 22rpx;
  color: #9a9c98;
  margin-top: 8rpx;
}

.coupon-right {
  flex: 1;
  padding: 24rpx;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.coupon-name {
  font-size: 26rpx;
  font-weight: 500;
  color: #2c2e2c;
  margin-bottom: 8rpx;
}

.coupon-date {
  font-size: 22rpx;
  color: #9a9c98;
  margin-bottom: 8rpx;
}

.coupon-status {
  font-size: 22rpx;
  color: #5b8c5a;
}

.coupon-status.used {
  color: #bcbeb9;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 120rpx 32rpx;
}

.empty-text {
  font-size: 28rpx;
  color: #9a9c98;
  margin-top: 16rpx;
}
</style>
