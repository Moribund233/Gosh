<template>
  <view class="product-page">
    <view class="status-bar" :style="{ height: statusBarHeight + 'px' }"></view>

    <view class="nav-bar">
      <view class="back-btn" @tap="goBack">
        <tui-icon name="arrowleft" size="40" color="#2c2e2c"></tui-icon>
      </view>
      <text class="nav-title">商品详情</text>
      <view class="nav-right">
        <tui-icon name="more" size="36" color="#555755"></tui-icon>
      </view>
    </view>

    <scroll-view class="content" scroll-y>
      <swiper class="product-gallery" :indicator-dots="false" :autoplay="false" :duration="300" circular>
        <swiper-item v-for="(img, idx) in product.images" :key="idx">
          <image class="gallery-img" :src="img" mode="aspectFill"></image>
        </swiper-item>
      </swiper>
      <view class="gallery-dots" v-if="product.images?.length > 1">
        <view class="dot" :class="{ active: idx === galleryIndex }" v-for="(_, idx) in product.images" :key="idx"></view>
      </view>

      <view class="product-meta">
        <text class="product-title">{{ product.name }}</text>
        <text class="product-sub">产地：{{ product.origin }} | 保质期：{{ product.shelf_life }} | 已售 {{ formatSales(product.sales) }}+</text>

        <view class="price-block">
          <text class="price-now"><text class="price-sym">¥</text>{{ formatPrice(product.price) }}</text>
          <text class="price-old" v-if="product.original_price > product.price">¥{{ formatPrice(product.original_price) }}</text>
          <text class="price-sales">月销 {{ formatSales(product.sales) }}+</text>
        </view>

        <view class="promo-tags" v-if="product.is_new || product.is_hot">
          <text class="promo-tag" v-if="product.is_new">新品</text>
          <text class="promo-tag" v-if="product.is_hot">热卖</text>
          <text class="promo-tag warm">产地直发</text>
        </view>
      </view>

      <view class="spec-section">
        <text class="spec-label">规格</text>
        <view class="spec-options">
          <view
            class="spec-opt"
            :class="{ selected: selectedSku === sku.id }"
            v-for="sku in product.skus"
            :key="sku.id"
            @tap="selectSku(sku.id)"
          >
            <text>{{ sku.name }}</text>
          </view>
        </view>
      </view>

      <view class="qty-section">
        <text class="qty-label">数量</text>
        <view class="qty-selector">
          <view class="qty-btn" @tap="decreaseQty">-</view>
          <text class="qty-num">{{ quantity }}</text>
          <view class="qty-btn" @tap="increaseQty">+</view>
        </view>
      </view>

      <view class="info-section">
        <text class="info-title">商品信息</text>
        <view class="info-grid">
          <view class="info-item" v-if="product.brand">
            <text class="info-key">品牌</text>
            <text class="info-val">{{ product.brand }}</text>
          </view>
          <view class="info-item" v-if="product.origin">
            <text class="info-key">产地</text>
            <text class="info-val">{{ product.origin }}</text>
          </view>
          <view class="info-item" v-if="product.net_content">
            <text class="info-key">净含量</text>
            <text class="info-val">{{ product.net_content }}</text>
          </view>
          <view class="info-item" v-if="product.shelf_life">
            <text class="info-key">保质期</text>
            <text class="info-val">{{ product.shelf_life }}</text>
          </view>
          <view class="info-item" v-if="product.storage_method">
            <text class="info-key">储存方式</text>
            <text class="info-val">{{ product.storage_method }}</text>
          </view>
        </view>
      </view>

      <view class="review-section">
        <view class="review-header">
          <text class="review-title">用户评价</text>
          <text class="review-more">查看全部 &gt;</text>
        </view>
        <view class="review-card" v-for="review in reviews" :key="review.id">
          <view class="review-user">
            <text class="review-name">{{ review.user_name }}</text>
            <text class="review-date">{{ review.created_at }}</text>
            <text class="review-rating">{{ '★'.repeat(review.rating) }}</text>
          </view>
          <text class="review-content">{{ review.content }}</text>
        </view>
      </view>

      <view class="recommend-section">
        <text class="recommend-title">你可能还喜欢</text>
        <scroll-view class="scroll-x" scroll-x show-scrollbar="false">
          <view class="recommend-card" v-for="item in recommends" :key="item.id" @tap="navigateToProduct(item.id)">
            <image class="recommend-img" :src="item.images?.[0]" mode="aspectFill"></image>
            <text class="recommend-name">{{ item.name }}</text>
            <text class="recommend-price">¥{{ formatPrice(item.price) }}</text>
          </view>
        </scroll-view>
      </view>
    </scroll-view>

    <view class="bottom-bar">
      <view class="icon-btn" @tap="goHome">
        <tui-icon name="home" size="40" color="#9a9c98"></tui-icon>
        <text class="icon-label">首页</text>
      </view>
      <view class="icon-btn" @tap="goCart">
        <tui-icon name="cart" size="40" color="#9a9c98"></tui-icon>
        <text class="icon-label">购物车</text>
      </view>
      <view class="spacer"></view>
      <view class="btn-cart" @tap="addToCart">加入购物车</view>
      <view class="btn-buy" @tap="buyNow">立即购买</view>
    </view>

    <view class="toast" v-if="toastVisible">{{ toastMessage }}</view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getProductDetail, getProducts, type Product } from '@/services/product'
import { getReviews, type Review } from '@/services/review'
import { formatPrice, formatSales } from '@/utils/format'

const statusBarHeight = ref(uni.getStorageSync('statusBarHeight') || 44)

const product = ref<Product>({} as Product)
const reviews = ref<Review[]>([])
const recommends = ref<Product[]>([])
const selectedSku = ref<number | null>(null)
const quantity = ref(1)
const galleryIndex = ref(0)
const toastVisible = ref(false)
const toastMessage = ref('')

const currentSku = computed(() => product.value.skus?.find(s => s.id === selectedSku.value))

onMounted(async () => {
  const pages = getCurrentPages()
  const page = pages[pages.length - 1] as any
  const id = page?.$page?.options?.id || page?.options?.id
  if (id) {
    await loadProduct(Number(id))
  }
})

async function loadProduct(id: number) {
  try {
    product.value = await getProductDetail(id)
    if (product.value.skus?.length > 0) {
      selectedSku.value = product.value.skus[0].id
    }
    await Promise.all([
      loadReviews(id),
      loadRecommends(),
    ])
  } catch {}
}

async function loadReviews(productId: number) {
  try {
    const res = await getReviews({ product_id: productId, size: 2 })
    reviews.value = res.list
  } catch {}
}

async function loadRecommends() {
  try {
    const res = await getProducts({ size: 6, is_featured: true })
    recommends.value = res.list.filter(p => p.id !== product.value.id).slice(0, 6)
  } catch {}
}

function selectSku(id: number) {
  selectedSku.value = id
}

function increaseQty() {
  const maxStock = currentSku.value?.stock || 99
  if (quantity.value < maxStock) quantity.value++
}

function decreaseQty() {
  if (quantity.value > 1) quantity.value--
}

function showToast(msg: string) {
  toastMessage.value = msg
  toastVisible.value = true
  setTimeout(() => { toastVisible.value = false }, 1500)
}

function addToCart() {
  if (!currentSku.value) return
  showToast('已加入购物车')
}

function buyNow() {
  if (!currentSku.value) {
    uni.showToast({ title: '请选择规格', icon: 'none' })
    return
  }
  uni.navigateTo({ url: '/pages/sub-package-order/checkout/checkout' })
}

function goBack() {
  uni.navigateBack()
}

function goHome() {
  uni.switchTab({ url: '/pages/index/index' })
}

function goCart() {
  const token = uni.getStorageSync('token')
  if (!token) {
    uni.navigateTo({ url: '/pages/login/login' })
    return
  }
  uni.switchTab({ url: '/pages/cart/cart' })
}

function navigateToProduct(id: number) {
  uni.redirectTo({ url: `/pages/sub-package-product/product/product?id=${id}` })
}
</script>

<style scoped>
.product-page {
  min-height: 100vh;
  background: #f7f8f5;
  display: flex;
  flex-direction: column;
}

.status-bar {
  background: #ffffff;
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
  border-radius: 50%;
  background: rgba(255,255,255,0.9);
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

.content {
  flex: 1;
}

.product-gallery {
  width: 100%;
  height: 750rpx;
  background: #f0f2ed;
}

.gallery-img {
  width: 100%;
  height: 100%;
}

.gallery-dots {
  position: absolute;
  bottom: 20rpx;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 8rpx;
}

.dot {
  width: 12rpx;
  height: 12rpx;
  border-radius: 6rpx;
  background: rgba(0,0,0,0.15);
  transition: all 0.2s;
}

.dot.active {
  width: 36rpx;
  background: #5b8c5a;
}

.product-meta {
  padding: 32rpx;
  background: #ffffff;
}

.product-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #2c2e2c;
  line-height: 1.4;
  display: block;
  margin-bottom: 12rpx;
}

.product-sub {
  font-size: 24rpx;
  color: #9a9c98;
  display: block;
  margin-bottom: 20rpx;
}

.price-block {
  display: flex;
  align-items: baseline;
  gap: 16rpx;
  margin-bottom: 16rpx;
}

.price-now {
  font-size: 48rpx;
  font-weight: 700;
  color: #d9534f;
}

.price-sym {
  font-size: 28rpx;
}

.price-old {
  font-size: 24rpx;
  color: #bcbeb9;
  text-decoration: line-through;
}

.price-sales {
  font-size: 24rpx;
  color: #9a9c98;
  margin-left: auto;
}

.promo-tags {
  display: flex;
  gap: 12rpx;
  flex-wrap: wrap;
}

.promo-tag {
  padding: 4rpx 16rpx;
  border-radius: 8rpx;
  font-size: 22rpx;
  background: #e8f0e6;
  color: #5b8c5a;
}

.promo-tag.warm {
  background: #f7efe0;
  color: #b88940;
}

.spec-section {
  padding: 32rpx;
  background: #ffffff;
  margin-top: 16rpx;
}

.spec-label {
  font-size: 24rpx;
  color: #555755;
  font-weight: 500;
  display: block;
  margin-bottom: 20rpx;
}

.spec-options {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.spec-opt {
  padding: 16rpx 28rpx;
  border: 2rpx solid #e5e7e3;
  border-radius: 999rpx;
  font-size: 24rpx;
  color: #555755;
  transition: all 0.15s;
}

.spec-opt.selected {
  border-color: #5b8c5a;
  background: #e8f0e6;
  color: #5b8c5a;
  font-weight: 500;
}

.qty-section {
  padding: 32rpx;
  background: #ffffff;
  margin-top: 2rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-top: 2rpx solid #eef0eb;
}

.qty-label {
  font-size: 24rpx;
  color: #555755;
  font-weight: 500;
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
  font-size: 32rpx;
  color: #2c2e2c;
  background: #ffffff;
}

.qty-btn:active {
  background: #f0f2ed;
}

.qty-num {
  width: 72rpx;
  text-align: center;
  font-size: 28rpx;
  font-weight: 500;
  border-left: 2rpx solid #e5e7e3;
  border-right: 2rpx solid #e5e7e3;
  line-height: 56rpx;
  background: #ffffff;
}

.info-section {
  padding: 32rpx;
  background: #ffffff;
  margin-top: 16rpx;
}

.info-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #2c2e2c;
  display: block;
  margin-bottom: 24rpx;
}

.info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16rpx;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.info-key {
  font-size: 22rpx;
  color: #9a9c98;
}

.info-val {
  font-size: 26rpx;
  color: #2c2e2c;
}

.review-section {
  padding: 32rpx;
  background: #ffffff;
  margin-top: 16rpx;
}

.review-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24rpx;
}

.review-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #2c2e2c;
}

.review-more {
  font-size: 24rpx;
  color: #9a9c98;
}

.review-card {
  background: #f0f2ed;
  padding: 20rpx;
  border-radius: 16rpx;
  margin-bottom: 16rpx;
}

.review-user {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-bottom: 12rpx;
}

.review-name {
  font-size: 24rpx;
  font-weight: 600;
  color: #2c2e2c;
}

.review-date {
  font-size: 22rpx;
  color: #9a9c98;
}

.review-rating {
  margin-left: auto;
  font-size: 22rpx;
  color: #f0ad4e;
}

.review-content {
  font-size: 24rpx;
  color: #555755;
  line-height: 1.5;
}

.recommend-section {
  padding: 32rpx;
  padding-bottom: 40rpx;
}

.recommend-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #2c2e2c;
  display: block;
  margin-bottom: 20rpx;
}

.scroll-x {
  display: flex;
  gap: 20rpx;
  overflow-x: auto;
  white-space: nowrap;
}

.scroll-x::-webkit-scrollbar {
  display: none;
}

.recommend-card {
  width: 260rpx;
  background: #ffffff;
  border-radius: 16rpx;
  overflow: hidden;
  flex-shrink: 0;
}

.recommend-img {
  width: 260rpx;
  height: 260rpx;
  background: #f0f2ed;
}

.recommend-name {
  display: block;
  padding: 12rpx 16rpx 4rpx;
  font-size: 24rpx;
  color: #2c2e2c;
  font-weight: 500;
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

.bottom-bar {
  height: 100rpx;
  display: flex;
  align-items: center;
  background: #ffffff;
  border-top: 2rpx solid #eef0eb;
  padding: 0 20rpx;
  gap: 16rpx;
  flex-shrink: 0;
  padding-bottom: env(safe-area-inset-bottom, 0);
}

.icon-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4rpx;
  min-width: 88rpx;
}

.icon-label {
  font-size: 20rpx;
  color: #9a9c98;
}

.spacer {
  flex: 1;
}

.btn-cart {
  padding: 20rpx 40rpx;
  background: #4d7a4c;
  color: #ffffff;
  border-radius: 999rpx;
  font-size: 28rpx;
  font-weight: 600;
}

.btn-cart:active {
  opacity: 0.85;
}

.btn-buy {
  padding: 20rpx 40rpx;
  background: #5b8c5a;
  color: #ffffff;
  border-radius: 999rpx;
  font-size: 28rpx;
  font-weight: 600;
}

.btn-buy:active {
  opacity: 0.85;
}

.toast {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: rgba(0,0,0,0.75);
  color: #ffffff;
  padding: 20rpx 40rpx;
  border-radius: 16rpx;
  font-size: 24rpx;
  z-index: 999;
  pointer-events: none;
}
</style>
