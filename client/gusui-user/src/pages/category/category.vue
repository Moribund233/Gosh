<template>
  <view class="category-page">
    <view class="nav-bar">
      <text class="nav-title">商品分类</text>
      <view class="nav-right" @tap="navigateToSearch">
        <tui-icon name="search" size="36" color="#555755"></tui-icon>
      </view>
    </view>

    <view class="category-layout">
      <scroll-view class="cat-sidebar" scroll-y :scroll-into-view="'cat-' + activeIndex">
        <view
          class="cat-item"
          :class="{ active: idx === activeIndex }"
          v-for="(cat, idx) in categoryTree"
          :key="cat.id"
          :id="'cat-' + idx"
          @tap="selectCategory(idx)"
        >
          <text>{{ cat.name }}</text>
        </view>
      </scroll-view>

      <scroll-view class="cat-content" scroll-y @scrolltolower="onLoadMore">
        <view class="current-cat">
          <text class="cat-title">{{ currentCat?.name }}</text>
        </view>

        <view class="sub-cat-list" v-if="currentCat?.children?.length">
          <view
            class="sub-cat-tag"
            :class="{ active: subActive === sub.id }"
            v-for="sub in currentCat.children"
            :key="sub.id"
            @tap="selectSubCategory(sub.id)"
          >
            <text>{{ sub.name }}</text>
          </view>
        </view>

        <view class="cat-banner" v-if="currentBanner">
          <image class="cat-banner-img" :src="currentBanner" mode="aspectFill"></image>
        </view>

        <view class="cat-product-grid">
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
    </view>
    <AppTabbar :current="1" />
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import AppTabbar from '@/components/AppTabbar.vue'
import { getCategoryTree, type Category } from '@/services/category'
import { getProducts, type Product } from '@/services/product'
import { formatPrice, formatSales } from '@/utils/format'

const categoryTree = ref<Category[]>([])
const activeIndex = ref(0)
const subActive = ref<number | null>(null)
const products = ref<Product[]>([])
const loading = ref(false)
const page = ref(1)

const currentCat = computed(() => categoryTree.value[activeIndex.value] || null)
const currentBanner = ref('')

onMounted(async () => {
  try {
    categoryTree.value = await getCategoryTree()
    if (categoryTree.value.length > 0) {
      loadProducts()
    }
  } catch {
    categoryTree.value = []
  }
})

function selectCategory(idx: number) {
  activeIndex.value = idx
  subActive.value = null
  page.value = 1
  loadProducts()
}

function selectSubCategory(id: number) {
  subActive.value = id
  page.value = 1
  loadProducts()
}

async function loadProducts() {
  loading.value = true
  try {
    const catId = subActive.value || currentCat.value?.id
    const res = await getProducts({ page: page.value, size: 10, category_id: catId })
    if (page.value === 1) {
      products.value = res.list
    } else {
      products.value = [...products.value, ...res.list]
    }
  } catch {
    products.value = []
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

function navigateToProduct(id: number) {
  uni.navigateTo({ url: `/pages/sub-package-product/product/product?id=${id}` })
}
</script>

<style scoped>
.category-page {
  height: 100vh;
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
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.category-layout {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.cat-sidebar {
  width: 170rpx;
  background: #f0f2ed;
  border-right: 1rpx solid #eef0eb;
  flex-shrink: 0;
}

.cat-item {
  padding: 28rpx 16rpx;
  text-align: center;
  font-size: 26rpx;
  color: #555755;
  border-left: 4rpx solid transparent;
  transition: all 0.15s;
}

.cat-item.active {
  background: #ffffff;
  color: #5b8c5a;
  font-weight: 600;
  border-left-color: #5b8c5a;
}

.cat-content {
  flex: 1;
  padding: 20rpx 24rpx;
}

.current-cat {
  margin-bottom: 20rpx;
}

.cat-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #2c2e2c;
}

.sub-cat-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12rpx;
  margin-bottom: 24rpx;
}

.sub-cat-tag {
  padding: 12rpx 24rpx;
  border: 2rpx solid #e5e7e3;
  border-radius: 999rpx;
  font-size: 24rpx;
  color: #555755;
}

.sub-cat-tag.active {
  background: #e8f0e6;
  border-color: #5b8c5a;
  color: #5b8c5a;
}

.cat-banner {
  height: 180rpx;
  border-radius: 16rpx;
  overflow: hidden;
  margin-bottom: 24rpx;
}

.cat-banner-img {
  width: 100%;
  height: 100%;
}

.cat-product-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16rpx;
}

.product-card {
  background: #ffffff;
  border-radius: 16rpx;
  overflow: hidden;
}

.product-img {
  width: 100%;
  height: 300rpx;
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
