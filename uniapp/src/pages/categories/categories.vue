<template>
  <view class="page">
    <u-search
      placeholder="搜索分类"
      shape="square"
      :show-action="false"
      bg-color="#FDFBF9"
      input-align="left"
      height="36"
    ></u-search>

    <view class="content">
      <scroll-view scroll-y class="left-menu">
        <view
          v-for="(cat, i) in categories"
          :key="cat"
          class="menu-item"
          :class="{ active: activeCat === i }"
          @tap="activeCat = i"
        >
          <text>{{ cat }}</text>
        </view>
      </scroll-view>

      <scroll-view scroll-y class="right-content">
        <view class="sub-title">{{ categories[activeCat] }}</view>
        <view class="sub-grid">
          <view
            v-for="j in 8"
            :key="j"
            class="sub-item"
            @tap="goCategory(activeCat, j)"
          >
            <image class="sub-img" src="/static/placeholder.png" mode="aspectFill"></image>
            <text class="sub-name">子分类{{ j }}</text>
          </view>
        </view>
      </scroll-view>
    </view>
  </view>
</template>

<script lang="ts">
export default {
  data() {
    return {
      activeCat: 0,
      categories: [
        '全部', '坚果炒货', '茗茶', '蜂蜜冲饮',
        '糕点零食', '干货调料', '礼品套装',
      ] as string[],
    }
  },
  methods: {
    goCategory(_catIdx: number, _subIdx: number) {
      uni.navigateTo({ url: '/pages/search/search' })
    },
  },
}
</script>

<style lang="scss" scoped>
.page {
  background: $uni-bg;
  min-height: 100vh;
  padding: 12px 16px 0;
}

.content {
  display: flex;
  height: calc(100vh - 180px);
  margin-top: 12px;
  margin-left: -16px;
  margin-right: -16px;
}

.left-menu {
  width: 88px;
  background: $uni-bg-white;
  border-radius: 0 $uni-border-radius-lg $uni-border-radius-lg 0;
  flex-shrink: 0;
}

.menu-item {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 48px;
  font-size: $uni-font-size-sm;
  color: $uni-text-color-grey;
  border-left: 3px solid transparent;
}
.menu-item.active {
  color: $uni-color-primary;
  border-left-color: $uni-color-primary;
  background: $uni-color-primary-light;
  font-weight: 600;
}

.right-content {
  flex: 1;
  padding: 0 16px 16px;
}

.sub-title {
  font-size: $uni-font-size-lg;
  font-weight: 600;
  color: $uni-text-color;
  margin-bottom: 16px;
  margin-top: 8px;
}

.sub-grid {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 16px 8px;
}

.sub-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
}

.sub-img {
  width: 64px;
  height: 64px;
  border-radius: $uni-border-radius-base;
  background-color: #F0EDE8;
}

.sub-name {
  font-size: $uni-font-size-xs;
  color: $uni-text-color;
}
</style>
