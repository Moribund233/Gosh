<template>
  <view class="home-page">
    <view class="nav-bar">
      <text class="nav-title">谷穗</text>
      <view class="nav-right" @tap="navigateToSearch">
        <tui-icon name="search" size="36" color="#555755"></tui-icon>
      </view>
    </view>

    <scroll-view class="content" scroll-y @scrolltolower="onLoadMore">
      <view class="search-bar" @tap="navigateToSearch">
        <tui-icon name="search" size="28" color="#9a9c98"></tui-icon>
        <text class="search-placeholder">搜索干杂货、产地、品牌…</text>
      </view>

      <swiper class="banner-swiper" :indicator-dots="false" :autoplay="true" :interval="3500" :duration="400" circular>
        <swiper-item v-for="item in banners" :key="item.id">
          <view class="banner-slide">
            <image class="banner-img" :src="item.image" mode="aspectFill"></image>
            <view class="banner-overlay">
              <text class="banner-title">{{ item.title }}</text>
            </view>
          </view>
        </swiper-item>
      </swiper>

      <scroll-view class="scroll-x cat-scroll" scroll-x show-scrollbar="false">
        <view class="cat-entry" v-for="(cat, idx) in categories" :key="idx" @tap="navigateToCategory">
          <view class="cat-icon">{{ cat.icon }}</view>
          <text class="cat-label">{{ cat.name }}</text>
        </view>
      </scroll-view>

      <view class="divider"></view>

      <view class="brand-story" v-if="brandStory" @tap="navigateToBrandStory">
        <image class="brand-img" :src="brandStory.image" mode="aspectFill"></image>
        <view class="brand-info">
          <text class="brand-title">{{ brandStory.title }}</text>
          <text class="brand-desc">{{ brandStory.content }}</text>
        </view>
      </view>

      <view class="divider"></view>

      <view class="section-header">
        <text class="section-title">限时秒杀</text>
        <view class="flash-timer" v-if="flashEndTime">
          <tui-icon name="clock" size="24" color="#d9534f"></tui-icon>
          <tui-countdown :time="flashCountdown" color="#d9534f" bcolor="#d9534f" :showDay="false"></tui-countdown>
        </view>
        <text class="section-more">更多 &gt;</text>
      </view>
      <scroll-view class="scroll-x product-scroll" scroll-x show-scrollbar="false">
        <view class="flash-card" v-for="item in flashSales" :key="item.id" @tap="navigateToProduct(item.product_id)">
          <image class="flash-img" :src="item.product_image" mode="aspectFill"></image>
          <text class="flash-name">{{ item.product_name }}</text>
          <view class="flash-price-row">
            <text class="flash-price">¥{{ formatPrice(item.flash_price) }}</text>
            <text class="flash-original">¥{{ formatPrice(item.original_price) }}</text>
          </view>
        </view>
      </scroll-view>

      <view class="divider"></view>

      <view class="section-header">
        <text class="section-title">为你推荐</text>
        <text class="section-more">更多 &gt;</text>
      </view>
      <view class="product-grid">
        <view class="product-card" v-for="item in products" :key="item.id" @tap="navigateToProduct(item.id)">
          <image class="product-img" :src="item.images?.[0]" mode="aspectFill"></image>
          <view class="product-info">
            <text class="product-name">{{ item.name }}</text>
            <text class="product-origin">{{ item.origin }}</text>
            <view class="product-price-row">
              <text class="product-price">¥{{ formatPrice(item.price) }}</text>
              <text class="product-sales">已售 {{ formatSales(item.sales) }}</text>
            </view>
          </view>
        </view>
      </view>

      <view class="load-more" v-if="loading">
        <text class="load-more-text">加载中...</text>
      </view>
    </scroll-view>
    <AppTabbar :current="0" />
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import AppTabbar from '@/components/AppTabbar.vue'
import { getBanners, type Banner, getBrandStory, type BrandStory } from '@/services/banner'
import { getFlashSales, type FlashSale } from '@/services/flash-sale'
import { getProducts, type Product } from '@/services/product'
import { getCategoryTree, type Category } from '@/services/category'
import { formatPrice, formatSales } from '@/utils/format'

const banners = ref<Banner[]>([])
const brandStory = ref<BrandStory | null>(null)
const flashSales = ref<FlashSale[]>([])
const products = ref<Product[]>([])
const categories = ref<{ name: string; icon: string }[]>([])
const loading = ref(false)
const page = ref(1)

const categoriesMap = [
  { name: '米面粮油', icon: '🌾' },
  { name: '豆类杂粮', icon: '🫘' },
  { name: '菌菇干货', icon: '🍄' },
  { name: '坚果干果', icon: '🥜' },
  { name: '调味干货', icon: '🌶' },
  { name: '花草茶', icon: '🍵' },
  { name: '蜜饯果脯', icon: '🍯' },
  { name: '海产干货', icon: '🦐' },
]

const flashEndTime = ref(0)
const flashCountdown = computed(() => {
  const remaining = flashEndTime.value - Date.now()
  return remaining > 0 ? Math.floor(remaining / 1000) : 0
})

onMounted(async () => {
  await Promise.all([
    loadBanners(),
    loadBrandStory(),
    loadFlashSales(),
    loadProducts(),
  ])
})

async function loadBanners() {
  try {
    banners.value = await getBanners()
  } catch {}
}

async function loadBrandStory() {
  try {
    brandStory.value = await getBrandStory()
  } catch {}
}

async function loadFlashSales() {
  try {
    const list = await getFlashSales()
    flashSales.value = list
    if (list.length > 0) {
      flashEndTime.value = new Date(list[0].end_at).getTime()
    }
  } catch {}
}

async function loadProducts() {
  loading.value = true
  try {
    const res = await getProducts({ page: page.value, size: 10, is_featured: true })
    if (page.value === 1) {
      products.value = res.list
    } else {
      products.value = [...products.value, ...res.list]
    }
  } catch {
  } finally {
    loading.value = false
  }
}

function onLoadMore() {
  page.value++
  loadProducts()
}

function navigateToSearch() {
  uni.navigateTo({ url: '/pages/search/search' })
}

function navigateToCategory() {
  uni.switchTab({ url: '/pages/category/category' })
}

function navigateToProduct(id: number) {
  uni.navigateTo({ url: `/pages/sub-package-product/product/product?id=${id}` })
}

function navigateToBrandStory() {}
</script>

<style scoped>
.home-page {
  min-height: 100vh;
  background: #f7f8f5;
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
  font-size: 36rpx;
  font-weight: 700;
  color: #5b8c5a;
  letter-spacing: 4rpx;
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
  height: calc(100vh - 88rpx);
}

.search-bar {
  display: flex;
  align-items: center;
  gap: 12rpx;
  background: #f0f2ed;
  border-radius: 999rpx;
  padding: 16rpx 24rpx;
  margin: 20rpx 32rpx;
}

.search-placeholder {
  font-size: 28rpx;
  color: #bcbeb9;
}

.banner-swiper {
  margin: 0 32rpx 20rpx;
  border-radius: 16rpx;
  overflow: hidden;
  height: 280rpx;
}

.banner-slide {
  position: relative;
  width: 100%;
  height: 100%;
}

.banner-img {
  width: 100%;
  height: 100%;
}

.banner-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(91,140,90,0.15), rgba(212,231,208,0.3));
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  padding: 32rpx;
}

.banner-title {
  font-size: 40rpx;
  font-weight: 700;
  color: #2c2e2c;
}

.scroll-x {
  display: flex;
  gap: 20rpx;
  overflow-x: auto;
  padding: 0 32rpx 24rpx;
  white-space: nowrap;
}

.scroll-x::-webkit-scrollbar {
  display: none;
}

.scroll-x > * {
  flex-shrink: 0;
}

.cat-entry {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12rpx;
  width: 128rpx;
}

.cat-icon {
  width: 96rpx;
  height: 96rpx;
  border-radius: 24rpx;
  background: #e8f0e6;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 44rpx;
}

.cat-label {
  font-size: 22rpx;
  color: #555755;
}

.divider {
  height: 16rpx;
  background: #f0f2ed;
}

.brand-story {
  margin: 20rpx 32rpx;
  border-radius: 16rpx;
  overflow: hidden;
  position: relative;
  height: 200rpx;
}

.brand-img {
  width: 100%;
  height: 100%;
}

.brand-info {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(91,140,90,0.2), transparent);
  padding: 32rpx;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.brand-title {
  font-size: 32rpx;
  font-weight: 700;
  color: #2c2e2c;
  margin-bottom: 8rpx;
}

.brand-desc {
  font-size: 24rpx;
  color: #555755;
}

.section-header {
  display: flex;
  align-items: center;
  padding: 24rpx 32rpx 20rpx;
  gap: 12rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #2c2e2c;
}

.flash-timer {
  display: flex;
  align-items: center;
  gap: 8rpx;
}

.section-more {
  font-size: 24rpx;
  color: #9a9c98;
  margin-left: auto;
}

.product-scroll {
  padding-bottom: 24rpx;
}

.flash-card {
  width: 280rpx;
  background: #ffffff;
  border-radius: 16rpx;
  overflow: hidden;
}

.flash-img {
  width: 280rpx;
  height: 280rpx;
  background: #e8f0e6;
}

.flash-name {
  display: block;
  padding: 12rpx 16rpx 4rpx;
  font-size: 26rpx;
  font-weight: 500;
  color: #2c2e2c;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.flash-price-row {
  padding: 0 16rpx 16rpx;
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.flash-price {
  font-size: 32rpx;
  font-weight: 700;
  color: #d9534f;
}

.flash-original {
  font-size: 22rpx;
  color: #9a9c98;
  text-decoration: line-through;
}

.product-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20rpx;
  padding: 0 32rpx 32rpx;
}

.product-card {
  background: #ffffff;
  border-radius: 16rpx;
  overflow: hidden;
}

.product-img {
  width: 100%;
  height: 340rpx;
  background: #f0f2ed;
}

.product-info {
  padding: 12rpx 16rpx 20rpx;
}

.product-name {
  font-size: 26rpx;
  font-weight: 500;
  color: #2c2e2c;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: 4rpx;
}

.product-origin {
  font-size: 22rpx;
  color: #9a9c98;
  margin-bottom: 12rpx;
  display: block;
}

.product-price-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.product-price {
  font-size: 32rpx;
  font-weight: 700;
  color: #d9534f;
}

.product-sales {
  font-size: 22rpx;
  color: #9a9c98;
}

.load-more {
  padding: 32rpx;
  text-align: center;
}

.load-more-text {
  font-size: 24rpx;
  color: #9a9c98;
}
</style>
