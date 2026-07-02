<template>
  <view class="points-page">
    <scroll-view class="content" scroll-y>
      <view class="points-header">
        <tui-icon name="wealth" size="48" color="#5b8c5a"></tui-icon>
        <text class="points-balance">{{ info.balance }}</text>
        <text class="points-label">当前积分</text>
      </view>

      <view class="points-list">
        <view class="points-item" v-for="record in info.records" :key="record.id">
          <view class="points-icon-box">
            <tui-icon :name="record.amount > 0 ? 'add' : 'reduce'" size="32" :color="record.amount > 0 ? '#5b8c5a' : '#d9534f'"></tui-icon>
          </view>
          <view class="points-info">
            <text class="points-remark">{{ record.remark }}</text>
            <text class="points-time">{{ record.created_at }}</text>
          </view>
          <text class="points-amount" :class="{ add: record.amount > 0, minus: record.amount < 0 }">
            {{ record.amount > 0 ? '+' : '' }}{{ record.amount }}
          </text>
        </view>
      </view>

      <view class="empty-state" v-if="info.records.length === 0 && !loading">
        <tui-icon name="wealth" size="64" color="#bcbeb9"></tui-icon>
        <text class="empty-text">暂无积分记录</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getPoints, type PointsInfo, type PointsRecord } from '@/services/points'

const loading = ref(false)
const info = reactive<PointsInfo>({
  balance: 0,
  records: [],
})

onMounted(async () => {
  loading.value = true
  try {
    const data = await getPoints()
    info.balance = data.balance
    info.records = data.records
  } catch {
    info.balance = 0
    info.records = []
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.points-page {
  min-height: 100vh;
  background: #f7f8f5;
}

.content { height: 100vh; }

.points-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60rpx 32rpx;
  background: linear-gradient(180deg, #e8f0e6, #f7f8f5);
}

.points-balance {
  font-size: 64rpx;
  font-weight: 700;
  color: #2c2e2c;
  margin-top: 16rpx;
}

.points-label {
  font-size: 24rpx;
  color: #9a9c98;
  margin-top: 8rpx;
}

.points-list {
  padding: 16rpx 32rpx;
}

.points-item {
  display: flex;
  align-items: center;
  gap: 20rpx;
  padding: 24rpx 0;
  border-bottom: 2rpx solid #eef0eb;
}

.points-icon-box {
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
  background: #f0f2ed;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.points-info {
  flex: 1;
}

.points-remark {
  font-size: 26rpx;
  color: #2c2e2c;
  display: block;
  margin-bottom: 4rpx;
}

.points-time {
  font-size: 22rpx;
  color: #9a9c98;
}

.points-amount {
  font-size: 30rpx;
  font-weight: 600;
}

.points-amount.add {
  color: #5b8c5a;
}

.points-amount.minus {
  color: #d9534f;
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
