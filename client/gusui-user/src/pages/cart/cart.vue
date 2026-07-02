<template>
  <view class="cart-page">
    <view class="nav-bar">
      <text class="nav-title">购物车</text>
      <view class="nav-right">
        <text class="nav-action" @tap="toggleEdit">{{ isEdit ? '完成' : '管理' }}</text>
      </view>
    </view>

    <scroll-view class="content" scroll-y v-if="items.length > 0">
      <view class="cart-list">
        <view class="cart-item" v-for="item in items" :key="item.id">
          <view class="check" :class="{ checked: item.selected }" @tap="toggleCheck(item.id)">
            <view class="check-inner"></view>
          </view>
          <image class="thumb" :src="item.product_image" mode="aspectFill" @tap="navigateToProduct(item.product_id)"></image>
          <view class="info">
            <text class="item-name">{{ item.product_name }}</text>
            <text class="item-spec">{{ item.sku_name }}</text>
            <view class="item-row">
              <text class="item-price">¥{{ formatPrice(item.price) }}</text>
              <view class="qty-selector" v-if="!isEdit">
                <view class="qty-btn" @tap="decrease(item)">-</view>
                <text class="qty-num">{{ item.quantity }}</text>
                <view class="qty-btn" @tap="increase(item)">+</view>
              </view>
              <view class="delete-btn" v-else @tap="removeItem(item.id)">删除</view>
            </view>
          </view>
        </view>
      </view>

      <view class="recommend-section">
        <text class="recommend-title">—— 为您推荐 ——</text>
        <scroll-view class="scroll-x" scroll-x show-scrollbar="false">
          <view class="recommend-card" v-for="item in recommends" :key="item.id" @tap="navigateToProduct(item.id)">
            <image class="recommend-img" :src="item.images?.[0]" mode="aspectFill"></image>
            <text class="recommend-name">{{ item.name }}</text>
            <text class="recommend-price">¥{{ formatPrice(item.price) }}</text>
          </view>
        </scroll-view>
      </view>
    </scroll-view>

    <view class="empty-state" v-else-if="!loading">
      <tui-icon name="cart" size="80" color="#bcbeb9"></tui-icon>
      <text class="empty-text">购物车是空的</text>
      <text class="empty-sub">去首页看看有什么好吃的吧</text>
      <view class="btn-gohome" @tap="goHome">去首页</view>
    </view>

    <view class="cart-summary" v-if="items.length > 0">
      <view class="all-check" @tap="toggleAll">
        <view class="check" :class="{ checked: allSelected }">
          <view class="check-inner"></view>
        </view>
        <text class="all-check-label">全选</text>
      </view>
      <view class="total-section">
        <text class="total-label">合计：</text>
        <text class="total-amount"><text class="total-sym">¥</text>{{ formatPrice(selectedAmount) }}</text>
      </view>
      <view class="btn-checkout" :class="{ disabled: selectedCount === 0 }" @tap="goCheckout">
        结算 ({{ selectedCount }})
      </view>
    </view>
    <AppTabbar :current="2" />
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import AppTabbar from '@/components/AppTabbar.vue'
import { getCartList, updateCartQuantity, toggleCartItem, toggleCartAll, removeCartItem, type CartItem } from '@/services/cart'
import { getProducts, type Product } from '@/services/product'
import { formatPrice } from '@/utils/format'

const items = ref<CartItem[]>([])
const recommends = ref<Product[]>([])
const loading = ref(false)
const isEdit = ref(false)

const selectedItems = computed(() => items.value.filter(i => i.selected))
const selectedCount = computed(() => selectedItems.value.length)
const selectedAmount = computed(() =>
  selectedItems.value.reduce((sum, i) => sum + i.price * i.quantity, 0)
)
const allSelected = computed(() => items.value.length > 0 && items.value.every(i => i.selected))

onMounted(async () => {
  await Promise.all([loadCart(), loadRecommends()])
})

async function loadCart() {
  loading.value = true
  try {
    items.value = await getCartList()
  } catch {
    items.value = []
  } finally {
    loading.value = false
  }
}

async function loadRecommends() {
  try {
    const res = await getProducts({ size: 6, is_featured: true })
    recommends.value = res.list.slice(0, 6)
  } catch {}
}

function toggleEdit() {
  isEdit.value = !isEdit.value
}

async function toggleCheck(id: number) {
  const item = items.value.find(i => i.id === id)
  if (!item) return
  item.selected = !item.selected
  try {
    await toggleCartItem(id, item.selected)
  } catch {}
}

async function toggleAll() {
  const newVal = !allSelected.value
  items.value.forEach(i => { i.selected = newVal })
  try {
    await toggleCartAll(newVal)
  } catch {}
}

async function increase(item: CartItem) {
  if (item.quantity >= item.stock) return
  item.quantity++
  try {
    await updateCartQuantity(item.id, item.quantity)
  } catch {}
}

async function decrease(item: CartItem) {
  if (item.quantity <= 1) return
  item.quantity--
  try {
    await updateCartQuantity(item.id, item.quantity)
  } catch {}
}

async function removeItem(id: number) {
  try {
    await removeCartItem(id)
    items.value = items.value.filter(i => i.id !== id)
  } catch {}
}

function goCheckout() {
  if (selectedCount.value === 0) return
  uni.navigateTo({ url: '/pages/sub-package-order/checkout/checkout' })
}

function goHome() {
  uni.switchTab({ url: '/pages/index/index' })
}

function navigateToProduct(id: number) {
  uni.navigateTo({ url: `/pages/sub-package-product/product/product?id=${id}` })
}
</script>

<style scoped>
.cart-page {
  min-height: 100vh;
  background: #f7f8f5;
  display: flex;
  flex-direction: column;
}

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
}

.nav-action {
  font-size: 26rpx;
  color: #555755;
}

.content { flex: 1; }

.cart-list {
  padding-top: 8rpx;
}

.cart-item {
  display: flex;
  gap: 20rpx;
  padding: 24rpx 32rpx;
  background: #ffffff;
  border-bottom: 2rpx solid #eef0eb;
  align-items: center;
}

.check {
  width: 40rpx;
  height: 40rpx;
  border-radius: 50%;
  border: 4rpx solid #e5e7e3;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
}

.check.checked {
  background: #5b8c5a;
  border-color: #5b8c5a;
}

.check-inner {
  width: 16rpx;
  height: 8rpx;
  border-left: 4rpx solid #fff;
  border-bottom: 4rpx solid #fff;
  transform: rotate(-45deg) translateY(-2rpx);
}

.thumb {
  width: 140rpx;
  height: 140rpx;
  border-radius: 16rpx;
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
  margin-bottom: 16rpx;
}

.item-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.item-price {
  font-size: 30rpx;
  font-weight: 700;
  color: #d9534f;
}

.qty-selector {
  display: flex;
  align-items: center;
  border: 2rpx solid #e5e7e3;
  border-radius: 8rpx;
  overflow: hidden;
}

.qty-btn {
  width: 56rpx;
  height: 56rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28rpx;
  color: #2c2e2c;
  background: #ffffff;
}

.qty-btn:active { background: #f0f2ed; }

.qty-num {
  width: 64rpx;
  text-align: center;
  font-size: 26rpx;
  font-weight: 500;
  border-left: 2rpx solid #e5e7e3;
  border-right: 2rpx solid #e5e7e3;
  line-height: 56rpx;
  background: #ffffff;
}

.delete-btn {
  padding: 12rpx 28rpx;
  background: #d9534f;
  color: #fff;
  border-radius: 8rpx;
  font-size: 24rpx;
}

.recommend-section {
  padding: 32rpx;
}

.recommend-title {
  display: block;
  text-align: center;
  font-size: 24rpx;
  color: #9a9c98;
  margin-bottom: 20rpx;
}

.scroll-x {
  display: flex;
  gap: 20rpx;
  overflow-x: auto;
  white-space: nowrap;
}

.scroll-x::-webkit-scrollbar { display: none; }

.recommend-card {
  width: 240rpx;
  background: #ffffff;
  border-radius: 16rpx;
  overflow: hidden;
  flex-shrink: 0;
}

.recommend-img {
  width: 240rpx;
  height: 240rpx;
  background: #f0f2ed;
}

.recommend-name {
  display: block;
  padding: 12rpx 16rpx 4rpx;
  font-size: 24rpx;
  color: #2c2e2c;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.recommend-price {
  display: block;
  padding: 0 16rpx 16rpx;
  font-size: 28rpx;
  font-weight: 700;
  color: #d9534f;
}

.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120rpx 32rpx;
}

.empty-text {
  font-size: 28rpx;
  color: #9a9c98;
  margin-top: 20rpx;
}

.empty-sub {
  font-size: 24rpx;
  color: #bcbeb9;
  margin-top: 8rpx;
}

.btn-gohome {
  margin-top: 32rpx;
  padding: 16rpx 48rpx;
  background: #5b8c5a;
  color: #fff;
  border-radius: 999rpx;
  font-size: 28rpx;
}

.cart-summary {
  height: 100rpx;
  display: flex;
  align-items: center;
  background: #ffffff;
  border-top: 2rpx solid #eef0eb;
  padding: 0 32rpx;
  gap: 20rpx;
  padding-bottom: env(safe-area-inset-bottom, 0);
}

.all-check {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.all-check .check {
  width: 36rpx;
  height: 36rpx;
}

.all-check-label {
  font-size: 24rpx;
  color: #555755;
}

.total-section {
  flex: 1;
  text-align: right;
  display: flex;
  align-items: baseline;
  justify-content: flex-end;
  gap: 4rpx;
}

.total-label {
  font-size: 24rpx;
  color: #555755;
}

.total-amount {
  font-size: 40rpx;
  font-weight: 700;
  color: #d9534f;
}

.total-sym {
  font-size: 28rpx;
}

.btn-checkout {
  padding: 20rpx 48rpx;
  background: #5b8c5a;
  color: #fff;
  border-radius: 999rpx;
  font-size: 28rpx;
  font-weight: 600;
  white-space: nowrap;
}

.btn-checkout.disabled {
  opacity: 0.4;
}
</style>
