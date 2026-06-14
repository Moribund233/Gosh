<template>
  <view class="page">
    <view class="user-card">
      <u-avatar size="56" bg-color="rgba(255,255,255,0.3)" fontSize="20"></u-avatar>
      <view class="user-info">
        <text class="user-name">未登录</text>
        <text class="user-hint">点击登录账号</text>
      </view>
      <u-icon name="arrow-right" color="rgba(255,255,255,0.6)" size="18"></u-icon>
    </view>

    <view class="stats-card">
      <view
        v-for="s in orderStats"
        :key="s.label"
        class="stat-item"
        @tap="goOrders(s.status)"
      >
        <text class="stat-num">{{ s.num }}</text>
        <text class="stat-label">{{ s.label }}</text>
      </view>
    </view>

    <u-cell-group :border="false">
      <u-cell
        v-for="m in menus"
        :key="m.label"
        :title="m.label"
        :icon="m.icon"
        icon-color="#8B5E3C"
        is-link
        :border="false"
        @click="onMenu(m.route)"
      ></u-cell>
    </u-cell-group>
  </view>
</template>

<script lang="ts">
export default {
  data() {
    return {
      orderStats: [
        { label: '待付款', status: 'unpaid', num: 0 },
        { label: '待发货', status: 'undelivered', num: 0 },
        { label: '待收货', status: 'delivering', num: 0 },
        { label: '待评价', status: 'rated', num: 0 },
      ] as OrderStat[],
      menus: [
        { label: '我的订单', icon: 'order', route: '/pages/order/list' },
        { label: '收货地址', icon: 'map', route: '/pages/address/list' },
        { label: '我的收藏', icon: 'star', route: '/pages/favorite/list' },
        { label: '联系客服', icon: 'chat', route: '' },
        { label: '设置', icon: 'setting', route: '/pages/setting/index' },
      ] as MenuItem[],
    }
  },
  methods: {
    goOrders(_status: string) {
      uni.navigateTo({ url: '/pages/order/list' })
    },
    onMenu(route: string) {
      if (route) uni.navigateTo({ url: route })
    },
  },
}

interface OrderStat {
  label: string
  status: string
  num: number
}

interface MenuItem {
  label: string
  icon: string
  route: string
}
</script>

<style lang="scss" scoped>
.page {
  background: $uni-bg;
  min-height: 100vh;
}

.user-card {
  display: flex;
  align-items: center;
  gap: 12px;
  background: linear-gradient(135deg, $uni-color-primary, #A0704E);
  padding: 32px 20px 24px;
}

.user-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.user-name {
  font-size: $uni-font-size-xl;
  font-weight: 600;
  color: #fff;
}

.user-hint {
  font-size: $uni-font-size-sm;
  color: rgba(255, 255, 255, 0.7);
}

.stats-card {
  display: flex;
  background: $uni-bg-white;
  border-radius: $uni-border-radius-lg;
  margin: -16px 16px 12px;
  padding: 16px 0;
  box-shadow: $uni-shadow-sm;
  position: relative;
  z-index: 1;
}

.stat-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.stat-num {
  font-size: $uni-font-size-xl;
  font-weight: 700;
  color: $uni-text-color;
}

.stat-label {
  font-size: $uni-font-size-xs;
  color: $uni-text-color-grey;
}

:deep(.u-cell) {
  background: $uni-bg-white;
  margin: 0 16px;
  padding: 4px 0;
}
:deep(.u-cell:first-child) {
  border-radius: $uni-border-radius-lg $uni-border-radius-lg 0 0;
}
:deep(.u-cell:last-child) {
  border-radius: 0 0 $uni-border-radius-lg $uni-border-radius-lg;
}
:deep(.u-cell-group) {
  background: transparent;
}
</style>
