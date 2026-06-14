<template>
  <view class="page">
    <u-search
      placeholder="搜索商品"
      shape="square"
      :show-action="false"
      bg-color="#FDFBF9"
      input-align="left"
      height="36"
      @search="onSearch"
    ></u-search>

    <u-scroll-list :indicator="false" class="cat-scroll">
      <view
        v-for="(cat, i) in categories"
        :key="cat"
        class="cat-chip"
        :class="{ active: activeCat === i }"
        @tap="activeCat = i"
      >
        {{ cat }}
      </view>
    </u-scroll-list>

    <u-gap height="8" bg-color="#F5F0EB"></u-gap>

    <view class="section-header">
      <text class="section-title">推荐商品</text>
      <text class="section-more">查看更多</text>
    </view>

    <view class="product-grid">
      <view
        v-for="(item, i) in products"
        :key="i"
        class="product-card"
        @tap="goDetail(i)"
      >
        <image class="product-img" :src="item.image" mode="aspectFill"></image>
        <view class="product-info">
          <text class="product-name">{{ item.name }}</text>
          <text class="product-price">
            <text class="price-symbol">¥</text>{{ item.price }}
          </text>
          <text class="product-sales">已售 {{ item.sales }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script lang="ts">
export default {
  data() {
    return {
      activeCat: 0,
      categories: ['全部', '坚果', '茶叶', '蜂蜜', '糕点', '干货', '礼盒'],
      products: [
        { name: '原味坚果礼盒 500g', price: '128', sales: 2345, image: '/static/placeholder.png' },
        { name: '有机正山小种红茶', price: '89', sales: 1876, image: '/static/placeholder.png' },
        { name: '百花蜜 500g 纯天然', price: '68', sales: 3210, image: '/static/placeholder.png' },
        { name: '手工蔓越莓曲奇 200g', price: '56', sales: 4567, image: '/static/placeholder.png' },
      ] as ProductItem[],
    }
  },
  onLoad() {
    this.loadProducts()
  },
  methods: {
    onSearch(val: string) {
      uni.navigateTo({ url: '/pages/search/search' })
    },
    goDetail(_i: number) {
      uni.navigateTo({ url: '/pages/product/detail' })
    },
    loadProducts() {
      // TODO: fetch from API
    },
  },
}

interface ProductItem {
  name: string
  price: string
  sales: number
  image: string
}
</script>

<style lang="scss" scoped>
.page {
  padding: 12px 16px;
  background: $uni-bg;
  min-height: 100vh;
}

.cat-scroll {
  margin-top: 12px;
}
.cat-chip {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 32px;
  padding: 0 16px;
  border-radius: $uni-border-radius-pill;
  background: $uni-bg-white;
  font-size: $uni-font-size-sm;
  color: $uni-text-color;
  margin-right: 10px;
  box-shadow: $uni-shadow-sm;
  white-space: nowrap;
}
.cat-chip.active {
  background: $uni-color-primary;
  color: #fff;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding: 0 4px;
}
.section-title {
  font-size: $uni-font-size-xl;
  font-weight: 600;
  color: $uni-text-color;
}
.section-more {
  font-size: $uni-font-size-sm;
  color: $uni-text-color-grey;
}

.product-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}
.product-card {
  background: $uni-bg-white;
  border-radius: $uni-border-radius-lg;
  overflow: hidden;
  box-shadow: $uni-shadow-sm;
}
.product-img {
  width: 100%;
  aspect-ratio: 1;
  background-color: #F0EDE8;
}
.product-info {
  padding: 8px 10px 12px;
}
.product-name {
  font-size: $uni-font-size-sm;
  font-weight: 500;
  color: $uni-text-color;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.4;
}
.product-price {
  display: block;
  font-size: $uni-font-size-lg;
  font-weight: 700;
  color: $uni-color-price;
  margin-top: 4px;
}
.price-symbol {
  font-size: $uni-font-size-xs;
  font-weight: 500;
}
.product-sales {
  font-size: $uni-font-size-xs;
  color: $uni-text-color-grey;
  margin-top: 2px;
}
</style>
