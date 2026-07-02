<template>
  <view class="address-page">
    <scroll-view class="content" scroll-y>
      <view class="address-list">
        <view class="address-card" v-for="item in list" :key="item.id">
          <view class="card-main" @tap="selectAddress(item)">
            <view class="info">
              <view class="name-row">
                <text class="name">{{ item.name }}</text>
                <text class="phone">{{ item.phone }}</text>
                <text class="default-tag" v-if="item.is_default">默认</text>
              </view>
              <text class="address-detail">{{ item.province }}{{ item.city }}{{ item.district }}{{ item.detail }}</text>
            </view>
          </view>
          <view class="card-actions">
            <view class="action-btn" @tap="setDefault(item.id)" v-if="!item.is_default">
              <tui-icon name="check" size="28" color="#9a9c98"></tui-icon>
              <text class="action-text">设为默认</text>
            </view>
            <view class="action-btn" @tap="editAddress(item.id)">
              <tui-icon name="edit" size="28" color="#9a9c98"></tui-icon>
              <text class="action-text">编辑</text>
            </view>
            <view class="action-btn" @tap="removeAddress(item.id)">
              <tui-icon name="delete" size="28" color="#d9534f"></tui-icon>
              <text class="action-text danger">删除</text>
            </view>
          </view>
        </view>
      </view>

      <view class="empty-state" v-if="list.length === 0 && !loading">
        <tui-icon name="location" size="64" color="#bcbeb9"></tui-icon>
        <text class="empty-text">暂无收货地址</text>
        <text class="empty-sub">添加地址以便下单</text>
      </view>

      <view class="add-btn-wrapper">
        <view class="btn-add" @tap="addAddress">
          <text class="btn-add-text">+ 新增收货地址</text>
        </view>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getAddressList, deleteAddress, setDefaultAddress, type Address } from '@/services/address'

const list = ref<Address[]>([])
const loading = ref(false)

onMounted(() => {
  loadList()
})

async function loadList() {
  loading.value = true
  try {
    list.value = await getAddressList()
  } catch {
    list.value = []
  } finally {
    loading.value = false
  }
}

function selectAddress(item: Address) {
  const pages = getCurrentPages()
  const prevPage = pages[pages.length - 2]
  if (prevPage) {
    (prevPage as any).$vm?.onAddressSelected?.(item)
    uni.navigateBack()
  }
}

function addAddress() {
  uni.navigateTo({ url: '/pages/address-edit/address-edit' })
}

function editAddress(id: number) {
  uni.navigateTo({ url: `/pages/address-edit/address-edit?id=${id}` })
}

async function removeAddress(id: number) {
  uni.showModal({
    title: '提示',
    content: '确定删除该地址吗？',
    success: async (res) => {
      if (res.confirm) {
        try {
          await deleteAddress(id)
          list.value = list.value.filter(i => i.id !== id)
        } catch {}
      }
    },
  })
}

async function setDefault(id: number) {
  try {
    await setDefaultAddress(id)
    list.value.forEach(i => { i.is_default = i.id === id })
  } catch {}
}
</script>

<style scoped>
.address-page {
  min-height: 100vh;
  background: #f7f8f5;
}

.content {
  height: 100vh;
}

.address-list {
  padding: 20rpx 32rpx;
}

.address-card {
  background: #ffffff;
  border-radius: 16rpx;
  margin-bottom: 20rpx;
  overflow: hidden;
}

.card-main {
  padding: 28rpx;
}

.info {
  flex: 1;
}

.name-row {
  display: flex;
  align-items: center;
  gap: 16rpx;
  margin-bottom: 8rpx;
}

.name {
  font-size: 28rpx;
  font-weight: 600;
  color: #2c2e2c;
}

.phone {
  font-size: 24rpx;
  color: #9a9c98;
}

.default-tag {
  font-size: 20rpx;
  color: #5b8c5a;
  background: #e8f0e6;
  padding: 2rpx 12rpx;
  border-radius: 8rpx;
}

.address-detail {
  font-size: 24rpx;
  color: #555755;
  line-height: 1.5;
}

.card-actions {
  display: flex;
  border-top: 2rpx solid #eef0eb;
}

.action-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8rpx;
  padding: 20rpx;
}

.action-btn:active {
  background: #f0f2ed;
}

.action-text {
  font-size: 24rpx;
  color: #555755;
}

.action-text.danger {
  color: #d9534f;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 120rpx 32rpx 60rpx;
}

.empty-text {
  font-size: 28rpx;
  color: #9a9c98;
  margin-top: 16rpx;
}

.empty-sub {
  font-size: 24rpx;
  color: #bcbeb9;
  margin-top: 8rpx;
}

.add-btn-wrapper {
  padding: 20rpx 32rpx 60rpx;
}

.btn-add {
  width: 100%;
  height: 88rpx;
  background: #5b8c5a;
  color: #ffffff;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-add:active {
  opacity: 0.85;
}

.btn-add-text {
  font-size: 28rpx;
  font-weight: 600;
}
</style>
