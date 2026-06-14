<template>
  <view class="page">
    <u-empty v-if="!items.length" text="购物车是空的" mode="car">
      <template #bottom>
        <u-button type="primary" shape="circle" @tap="goHome">去逛逛</u-button>
      </template>
    </u-empty>

    <view v-for="(item, i) in items" :key="item.id" class="cart-item">
      <u-checkbox
        v-model="item.checked"
        shape="circle"
        active-color="#8B5E3C"
      ></u-checkbox>
      <image class="item-img" :src="item.image" mode="aspectFill"></image>
      <view class="item-info">
        <text class="item-name">{{ item.name }}</text>
        <text class="item-spec">{{ item.spec }}</text>
        <view class="item-bottom">
          <text class="item-price">¥{{ item.price }}</text>
          <u-number-box
            v-model="item.qty"
            :min="1"
            :max="99"
            :step="1"
            size="24"
            input-width="32"
            bg-color="#FDFBF9"
          ></u-number-box>
        </view>
      </view>
    </view>

    <u-gap height="100" bg-color="#F5F0EB"></u-gap>

    <view class="cart-bar">
      <u-checkbox
        :model-value="allChecked"
        shape="circle"
        active-color="#8B5E3C"
        label="全选"
        label-size="14"
        @change="toggleAll"
      ></u-checkbox>
      <view class="bar-right">
        <text class="total-label">合计: </text>
        <text class="total-price">¥{{ total }}</text>
        <u-button type="primary" shape="circle" :disabled="!total" @tap="checkout">结算</u-button>
      </view>
    </view>
  </view>
</template>

<script lang="ts">
export default {
  data() {
    return {
      items: [
        { id: 1, name: '原味坚果礼盒 500g', spec: '500g/盒', price: 128, qty: 1, checked: true, image: '/static/placeholder.png' },
        { id: 2, name: '有机正山小种红茶', spec: '250g/罐', price: 89, qty: 2, checked: false, image: '/static/placeholder.png' },
      ] as CartItem[],
    }
  },
  computed: {
    allChecked(): boolean {
      return this.items.length > 0 && this.items.every((i) => i.checked)
    },
    total(): string {
      return this.items
        .filter((i) => i.checked)
        .reduce((s, i) => s + i.price * i.qty, 0)
        .toFixed(2)
    },
  },
  methods: {
    toggleAll(checked: boolean) {
      this.items.forEach((i) => (i.checked = checked))
    },
    goHome() {
      uni.switchTab({ url: '/pages/index/index' })
    },
    checkout() {
      uni.navigateTo({ url: '/pages/order/checkout' })
    },
  },
}

interface CartItem {
  id: number
  name: string
  spec: string
  price: number
  qty: number
  checked: boolean
  image: string
}
</script>

<style lang="scss" scoped>
.page {
  background: $uni-bg;
  min-height: 100vh;
  padding: 12px 16px;
}

.cart-item {
  display: flex;
  align-items: center;
  gap: 10px;
  background: $uni-bg-white;
  border-radius: $uni-border-radius-lg;
  padding: 12px;
  margin-bottom: 8px;
  box-shadow: $uni-shadow-sm;
}

.item-img {
  width: 80px;
  height: 80px;
  border-radius: $uni-border-radius-base;
  background-color: #F0EDE8;
  flex-shrink: 0;
}

.item-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.item-name {
  font-size: $uni-font-size-sm;
  font-weight: 500;
  color: $uni-text-color;
}

.item-spec {
  font-size: $uni-font-size-xs;
  color: $uni-text-color-grey;
}

.item-bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 4px;
}

.item-price {
  font-size: $uni-font-size-lg;
  font-weight: 700;
  color: $uni-color-price;
}

.cart-bar {
  position: fixed;
  bottom: 72px;
  left: 0;
  right: 0;
  background: $uni-bg-white;
  border-top: 1px solid $uni-border-color;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
}

.bar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.total-label {
  font-size: $uni-font-size-sm;
  color: $uni-text-color;
}

.total-price {
  font-size: $uni-font-size-xl;
  font-weight: 700;
  color: $uni-color-price;
  margin-right: 4px;
}
</style>
