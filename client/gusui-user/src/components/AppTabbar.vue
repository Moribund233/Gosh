<template>
  <tui-tabbar
    :current="current"
    :tabBar="tabList"
    color="#9a9c98"
    selectedColor="#5b8c5a"
    backgroundColor="#ffffff"
    @click="onTabClick"
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

const props = defineProps<{
  current: number
}>()

const authStore = useAuthStore()

const cartBadge = computed(() => 0)

const tabList = computed(() => [
  { name: 'home', activeName: 'home-fill', text: '首页', pagePath: '/pages/index/index' },
  { name: 'category', activeName: 'category-fill', text: '分类', pagePath: '/pages/category/category' },
  { name: 'cart', activeName: 'cart-fill', text: '购物车', pagePath: '/pages/cart/cart', num: cartBadge.value },
  { name: 'people', activeName: 'people-fill', text: '我的', pagePath: '/pages/profile/profile' },
])

function onTabClick(e: any) {
  const { pagePath } = e
  if (pagePath === '/pages/cart/cart' && !authStore.isLoggedIn) {
    uni.navigateTo({ url: '/pages/login/login' })
    return
  }
  uni.reLaunch({ url: pagePath })
}
</script>

<style>
.tui-tabbar-item {
  justify-content: center !important;
  gap: 2rpx;
}
.tui-tabbar-item .tui-text-scale {
  line-height: 26rpx !important;
}
</style>
