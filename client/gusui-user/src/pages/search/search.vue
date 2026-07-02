<template>
  <view class="search-page">
    <view class="status-bar" :style="{ height: statusBarHeight + 'px' }"></view>

    <view class="search-header">
      <view class="search-input-bar">
        <tui-icon name="search" size="28" color="#5b8c5a"></tui-icon>
        <input
          class="search-input"
          v-model="keyword"
          type="text"
          placeholder="搜索干杂货、产地、品牌…"
          placeholder-class="input-placeholder"
          :focus="autoFocus"
          @confirm="onSearch"
          @input="onInputChange"
        />
        <view class="search-clear" v-if="keyword" @tap="clearKeyword">
          <tui-icon name="close-fill" size="28" color="#bcbeb9"></tui-icon>
        </view>
      </view>
      <text class="search-cancel" @tap="goBack">取消</text>
    </view>

    <view class="content" v-if="!showResults">
      <view class="section" v-if="hotTags.length > 0">
        <text class="section-label">热门搜索</text>
        <view class="tag-list">
          <text class="tag" v-for="(tag, idx) in hotTags" :key="idx" @tap="searchKeyword(tag)">{{ tag }}</text>
        </view>
      </view>

      <view class="divider"></view>

      <view class="section" v-if="historyList.length > 0">
        <view class="section-header-row">
          <text class="section-label">搜索历史</text>
          <text class="clear-btn" @tap="clearHistory">清空</text>
        </view>
        <view class="tag-list">
          <text class="tag history-tag" v-for="(item, idx) in historyList" :key="idx" @tap="searchKeyword(item)">{{ item }}</text>
        </view>
      </view>
    </view>

    <scroll-view class="content results" scroll-y v-else @scrolltolower="onLoadMore">
      <view class="result-item" v-for="item in results" :key="item.id" @tap="navigateToProduct(item.id)">
        <image class="result-img" :src="item.images?.[0]" mode="aspectFill"></image>
        <view class="result-info">
          <text class="result-name">{{ item.name }}</text>
          <text class="result-meta">产地：{{ item.origin }} | 月销 {{ formatSales(item.sales) }}+</text>
          <view class="result-price-row">
            <text class="result-price">¥{{ formatPrice(item.price) }}</text>
            <text class="result-original" v-if="item.original_price > item.price">¥{{ formatPrice(item.original_price) }}</text>
          </view>
        </view>
      </view>

      <view class="load-more" v-if="loading">
        <text class="load-more-text">加载中...</text>
      </view>

      <view class="empty-state" v-if="!loading && results.length === 0">
        <tui-icon name="search" size="64" color="#bcbeb9"></tui-icon>
        <text class="empty-text">没有找到相关商品</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { searchProducts, getHotSearch, getSearchHistory, clearSearchHistory, type Product } from '@/services/product'
import { formatPrice, formatSales } from '@/utils/format'

const statusBarHeight = ref(uni.getStorageSync('statusBarHeight') || 44)

const keyword = ref('')
const autoFocus = ref(false)
const showResults = ref(false)
const hotTags = ref<string[]>([])
const historyList = ref<string[]>([])
const results = ref<Product[]>([])
const loading = ref(false)
const page = ref(1)

onMounted(async () => {
  autoFocus.value = true
  await Promise.all([
    loadHotTags(),
    loadHistory(),
  ])
})

async function loadHotTags() {
  try {
    hotTags.value = await getHotSearch()
  } catch {
    hotTags.value = ['五常大米', '香菇', '黑木耳', '小米', '红枣', '枸杞']
  }
}

async function loadHistory() {
  try {
    historyList.value = await getSearchHistory()
  } catch {
    historyList.value = []
  }
}

function onInputChange(e: any) {
  const val = e.detail.value
  if (!val) {
    showResults.value = false
    results.value = []
  }
}

function searchKeyword(kw: string) {
  keyword.value = kw
  onSearch()
}

async function onSearch() {
  if (!keyword.value.trim()) return
  showResults.value = true
  page.value = 1
  loading.value = true
  try {
    const res = await searchProducts({ keyword: keyword.value, page: page.value, size: 20 })
    results.value = res.list
  } catch {
    results.value = []
  } finally {
    loading.value = false
  }
}

function clearKeyword() {
  keyword.value = ''
  showResults.value = false
  results.value = []
}

async function clearHistory() {
  try {
    await clearSearchHistory()
  } catch {}
  historyList.value = []
}

async function onLoadMore() {
  if (loading.value) return
  page.value++
  loading.value = true
  try {
    const res = await searchProducts({ keyword: keyword.value, page: page.value, size: 20 })
    results.value = [...results.value, ...res.list]
  } catch {
  } finally {
    loading.value = false
  }
}

function navigateToProduct(id: number) {
  uni.navigateTo({ url: `/pages/sub-package-product/product/product?id=${id}` })
}

function goBack() {
  uni.navigateBack()
}
</script>

<style scoped>
.search-page {
  min-height: 100vh;
  background: #f7f8f5;
}

.status-bar {
  background: #ffffff;
}

.search-header {
  display: flex;
  align-items: center;
  padding: 16rpx 32rpx;
  background: #ffffff;
  gap: 20rpx;
  border-bottom: 1rpx solid #eef0eb;
}

.search-input-bar {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 12rpx;
  background: #f0f2ed;
  border: 2rpx solid #5b8c5a;
  border-radius: 999rpx;
  padding: 12rpx 24rpx;
}

.search-input {
  flex: 1;
  font-size: 28rpx;
  color: #2c2e2c;
  background: transparent;
  border: none;
  outline: none;
  height: 48rpx;
}

.search-clear {
  width: 48rpx;
  height: 48rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.search-cancel {
  font-size: 28rpx;
  color: #555755;
  flex-shrink: 0;
}

.content {
  padding: 20rpx 0;
}

.content.results {
  padding: 0;
}

.section {
  padding: 0 32rpx 20rpx;
}

.section-label {
  font-size: 26rpx;
  color: #555755;
  font-weight: 600;
  display: block;
  margin-bottom: 20rpx;
}

.section-header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
}

.clear-btn {
  font-size: 24rpx;
  color: #9a9c98;
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.tag {
  padding: 12rpx 28rpx;
  border: 2rpx solid #e5e7e3;
  border-radius: 999rpx;
  font-size: 24rpx;
  color: #555755;
  background: #ffffff;
}

.tag:active {
  background: #e8f0e6;
  border-color: #5b8c5a;
  color: #5b8c5a;
}

.history-tag {
  background: #f0f2ed;
  border-color: transparent;
}

.divider {
  height: 16rpx;
  background: #f0f2ed;
}

.result-item {
  display: flex;
  gap: 20rpx;
  padding: 24rpx 32rpx;
  border-bottom: 2rpx solid #eef0eb;
  background: #ffffff;
}

.result-img {
  width: 160rpx;
  height: 160rpx;
  border-radius: 16rpx;
  background: #f0f2ed;
  flex-shrink: 0;
}

.result-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.result-name {
  font-size: 28rpx;
  font-weight: 500;
  color: #2c2e2c;
  margin-bottom: 8rpx;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.result-meta {
  font-size: 22rpx;
  color: #9a9c98;
  margin-bottom: 12rpx;
}

.result-price-row {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.result-price {
  font-size: 32rpx;
  font-weight: 700;
  color: #d9534f;
}

.result-original {
  font-size: 22rpx;
  color: #9a9c98;
  text-decoration: line-through;
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
  color: #9a9c98;
}

.empty-text {
  font-size: 28rpx;
  margin-top: 20rpx;
}

.input-placeholder {
  color: #bcbeb9;
  font-size: 28rpx;
}
</style>
