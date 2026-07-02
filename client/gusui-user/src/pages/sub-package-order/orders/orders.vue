<template>
  <view class="orders-page">
    <view class="status-bar" :style="{ height: statusBarHeight + 'px' }"></view>

    <view class="nav-bar">
      <view class="back-btn" @tap="goBack">
        <tui-icon name="arrowleft" size="40" color="#2c2e2c"></tui-icon>
      </view>
      <text class="nav-title">我的订单</text>
    </view>

    <view class="order-tabs">
      <view class="tab" :class="{ active: activeTab === tab.value }" v-for="tab in tabs" :key="tab.value" @tap="switchTab(tab.value)">
        <text>{{ tab.label }}</text>
      </view>
    </view>

    <scroll-view class="content" scroll-y @scrolltolower="onLoadMore">
      <view class="order-card" v-for="order in orders" :key="order.id">
        <view class="order-hd">
          <text class="store">谷穗优选</text>
          <text class="status" :class="getStatusClass(order.status)">{{ getStatusText(order.status) }}</text>
        </view>
        <view class="order-body">
          <image class="thumb" :src="order.items?.[0]?.product_image" mode="aspectFill"></image>
          <view class="info">
            <text class="item-name">{{ order.items?.[0]?.product_name || '商品' }}</text>
            <text class="item-spec">{{ order.items?.[0]?.sku_name || '' }}</text>
            <text class="order-time">下单时间：{{ order.created_at }}</text>
          </view>
          <view class="amount">
            <text class="order-total">¥{{ formatPrice(order.total_amount) }}</text>
          </view>
        </view>
        <view class="order-ft">
          <view class="btn-sm" @tap="handleCancel(order)" v-if="order.status === 'unpaid'">取消订单</view>
          <view class="btn-sm" @tap="handlePay(order)" v-if="order.status === 'unpaid'">立即付款</view>
          <view class="btn-sm" @tap="handleLogistics(order)" v-if="order.status === 'delivering'">查看物流</view>
          <view class="btn-sm primary" @tap="handleConfirm(order)" v-if="order.status === 'delivering'">确认收货</view>
          <view class="btn-sm" @tap="handleLogistics(order)" v-if="order.status === 'delivered'">查看物流</view>
          <view class="btn-sm primary" @tap="handleRebuy(order)" v-if="order.status === 'delivered'">再次购买</view>
        </view>
      </view>

      <view class="load-more" v-if="loading">
        <text class="load-more-text">加载中...</text>
      </view>

      <view class="empty-state" v-if="!loading && orders.length === 0">
        <tui-icon name="order" size="64" color="#bcbeb9"></tui-icon>
        <text class="empty-text">暂无订单</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getOrderList, cancelOrder, confirmOrder, type Order } from '@/services/order'
import { formatPrice } from '@/utils/format'

const statusBarHeight = ref(uni.getStorageSync('statusBarHeight') || 44)

const tabs = [
  { label: '全部', value: '' },
  { label: '待付款', value: 'unpaid' },
  { label: '待发货', value: 'undelivered' },
  { label: '待收货', value: 'delivering' },
  { label: '已完成', value: 'delivered' },
]

const activeTab = ref('')
const orders = ref<Order[]>([])
const loading = ref(false)
const page = ref(1)

onMounted(() => {
  const pages = getCurrentPages()
  const pageObj = pages[pages.length - 1] as any
  const status = pageObj?.$page?.options?.status || pageObj?.options?.status || ''
  if (status) activeTab.value = status
  loadOrders()
})

function switchTab(tab: string) {
  activeTab.value = tab
  page.value = 1
  loadOrders()
}

async function loadOrders() {
  loading.value = true
  try {
    const res = await getOrderList({ status: activeTab.value, page: page.value, size: 10 })
    if (page.value === 1) {
      orders.value = res.list
    } else {
      orders.value = [...orders.value, ...res.list]
    }
  } catch {
    orders.value = []
  } finally {
    loading.value = false
  }
}

function onLoadMore() {
  page.value++
  loadOrders()
}

function getStatusText(status: string): string {
  const map: Record<string, string> = {
    unpaid: '待付款',
    undelivered: '待发货',
    delivering: '待收货',
    delivered: '已完成',
    cancelled: '已取消',
    refund: '售后中',
  }
  return map[status] || status
}

function getStatusClass(status: string): string {
  if (status === 'delivered' || status === 'cancelled') return 'muted'
  return ''
}

function goBack() { uni.navigateBack() }

async function handleCancel(order: Order) {
  uni.showModal({
    title: '提示',
    content: '确定取消订单吗？',
    success: async (res) => {
      if (res.confirm) {
        try {
          await cancelOrder(order.id)
          uni.showToast({ title: '已取消', icon: 'success' })
          loadOrders()
        } catch {}
      }
    },
  })
}

function handlePay(_order: Order) {
  uni.showToast({ title: '支付功能开发中', icon: 'none' })
}

function handleConfirm(order: Order) {
  uni.showModal({
    title: '提示',
    content: '确认收货吗？',
    success: async (res) => {
      if (res.confirm) {
        try {
          await confirmOrder(order.id)
          uni.showToast({ title: '已确认收货', icon: 'success' })
          loadOrders()
        } catch {}
      }
    },
  })
}

function handleLogistics(_order: Order) {
  uni.showToast({ title: '物流查询开发中', icon: 'none' })
}

function handleRebuy(order: Order) {
  if (order.items?.[0]) {
    uni.navigateTo({ url: `/pages/sub-package-product/product/product?id=${order.items[0].product_id}` })
  }
}
</script>

<style scoped>
.orders-page {
  min-height: 100vh;
  background: #f7f8f5;
  display: flex;
  flex-direction: column;
}

.status-bar { background: #ffffff; }

.nav-bar {
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #ffffff;
  padding: 0 32rpx;
  position: relative;
  border-bottom: 1rpx solid #eef0eb;
}

.back-btn {
  position: absolute;
  left: 16rpx;
  top: 50%;
  transform: translateY(-50%);
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.nav-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #2c2e2c;
}

.order-tabs {
  display: flex;
  background: #ffffff;
  border-bottom: 1rpx solid #eef0eb;
  position: sticky;
  top: 0;
  z-index: 2;
}

.tab {
  flex: 1;
  text-align: center;
  padding: 24rpx 0;
  font-size: 26rpx;
  color: #555755;
  border-bottom: 4rpx solid transparent;
}

.tab.active {
  color: #5b8c5a;
  font-weight: 600;
  border-bottom-color: #5b8c5a;
}

.content { flex: 1; }

.order-card {
  background: #ffffff;
  margin: 20rpx 32rpx;
  border-radius: 16rpx;
  overflow: hidden;
}

.order-hd {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24rpx 24rpx 0;
  font-size: 24rpx;
}

.store {
  font-weight: 500;
  color: #2c2e2c;
}

.status {
  color: #5b8c5a;
  font-weight: 500;
}

.status.muted {
  color: #9a9c98;
}

.order-body {
  display: flex;
  gap: 20rpx;
  padding: 20rpx 24rpx;
  align-items: center;
}

.thumb {
  width: 128rpx;
  height: 128rpx;
  border-radius: 12rpx;
  background: #f0f2ed;
  flex-shrink: 0;
}

.info {
  flex: 1;
  min-width: 0;
}

.item-name {
  font-size: 26rpx;
  font-weight: 500;
  color: #2c2e2c;
  display: block;
  margin-bottom: 4rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.item-spec {
  font-size: 22rpx;
  color: #9a9c98;
  display: block;
  margin-bottom: 8rpx;
}

.order-time {
  font-size: 22rpx;
  color: #bcbeb9;
}

.amount {
  text-align: right;
  flex-shrink: 0;
}

.order-total {
  font-weight: 600;
  font-size: 28rpx;
  color: #d9534f;
}

.order-ft {
  display: flex;
  justify-content: flex-end;
  gap: 16rpx;
  padding: 16rpx 24rpx 24rpx;
  border-top: 2rpx solid #eef0eb;
}

.btn-sm {
  padding: 12rpx 28rpx;
  border: 2rpx solid #e5e7e3;
  border-radius: 999rpx;
  font-size: 24rpx;
  color: #555755;
  background: transparent;
}

.btn-sm:active {
  background: #f0f2ed;
}

.btn-sm.primary {
  background: #5b8c5a;
  color: #ffffff;
  border-color: #5b8c5a;
}

.load-more {
  padding: 32rpx;
  text-align: center;
}

.load-more-text {
  font-size: 24rpx;
  color: #9a9c98;
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
