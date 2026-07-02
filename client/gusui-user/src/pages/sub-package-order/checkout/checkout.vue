<template>
  <view class="checkout-page">
    <scroll-view class="content" scroll-y>
      <view class="section">
        <view class="address-card" @tap="selectAddress">
          <view class="address-icon">
            <tui-icon name="location" size="36" color="#5b8c5a"></tui-icon>
          </view>
          <view class="address-info" v-if="address">
            <view class="address-name-row">
              <text class="address-name">{{ address.name }}</text>
              <text class="address-phone">{{ address.phone }}</text>
              <text class="address-change">更换 &gt;</text>
            </view>
            <text class="address-detail">{{ address.province }}{{ address.city }}{{ address.district }}{{ address.detail }}</text>
          </view>
          <view class="address-info" v-else>
            <text class="address-empty">请选择收货地址</text>
          </view>
        </view>
      </view>

      <view class="section">
        <view class="sec-title">
          <tui-icon name="cart" size="28" color="#5b8c5a"></tui-icon>
          <text class="sec-title-text">商品清单 · 共{{ items.length }}件</text>
        </view>
        <view class="checkout-item" v-for="item in items" :key="item.sku_id">
          <image class="item-thumb" :src="item.product_image" mode="aspectFill"></image>
          <view class="item-info">
            <text class="item-name">{{ item.product_name }}</text>
            <text class="item-spec">{{ item.sku_name }}</text>
          </view>
          <text class="item-qty">×{{ item.quantity }}</text>
          <text class="item-price">¥{{ formatPrice(item.price * item.quantity) }}</text>
        </view>
      </view>

      <view class="section">
        <view class="coupon-row" @tap="selectCoupon">
          <text class="coupon-label">优惠券</text>
          <text class="coupon-value" :class="{ highlight: selectedCoupon }">{{ selectedCoupon ? '-¥' + formatPrice(selectedCoupon) : '选择优惠券' }}</text>
          <tui-icon name="arrowright" size="28" color="#bcbeb9"></tui-icon>
        </view>
      </view>

      <view class="section">
        <view class="fee-row">
          <text class="fee-label">商品小计</text>
          <text class="fee-value">¥{{ formatPrice(subtotal) }}</text>
        </view>
        <view class="fee-row">
          <text class="fee-label">运费</text>
          <text class="fee-value free">免运费</text>
        </view>
        <view class="fee-row" v-if="discount > 0">
          <text class="fee-label">优惠</text>
          <text class="fee-value discount">-¥{{ formatPrice(discount) }}</text>
        </view>
      </view>

      <view class="section">
        <view class="sec-title">
          <tui-icon name="wallet" size="28" color="#5b8c5a"></tui-icon>
          <text class="sec-title-text">支付方式</text>
        </view>
        <view class="payment-opt" @tap="selectPayment('wechat')">
          <view class="radio" :class="{ selected: payment === 'wechat' }">
            <view class="radio-inner"></view>
          </view>
          <tui-icon name="wechat" size="36" color="#5b8c5a"></tui-icon>
          <text class="pay-name">微信支付</text>
        </view>
      </view>

      <view class="section">
        <view class="remark-row">
          <text class="remark-label">订单备注</text>
          <input class="remark-input" v-model="remark" type="text" placeholder="选填，如有特殊要求请注明" placeholder-class="placeholder" />
        </view>
      </view>
    </scroll-view>

    <view class="bottom-bar">
      <view class="total-section">
        <text class="total-label">实付金额</text>
        <text class="total-amount"><text class="total-sym">¥</text>{{ formatPrice(subtotal - discount) }}</text>
      </view>
      <view class="btn-submit" :class="{ disabled: submitting }" @tap="submitOrder">
        <text v-if="submitting">提交中...</text>
        <text v-else>提交订单</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { createOrder } from '@/services/order'
import { getAddressList, type Address } from '@/services/address'
import { formatPrice } from '@/utils/format'

interface CheckoutItem {
  sku_id: number
  product_id: number
  product_name: string
  product_image: string
  sku_name: string
  price: number
  quantity: number
}

const address = ref<Address | null>(null)
const items = ref<CheckoutItem[]>([])
const payment = ref('wechat')
const remark = ref('')
const selectedCoupon = ref<number | null>(null)
const discount = ref(0)
const submitting = ref(false)

const subtotal = computed(() =>
  items.value.reduce((sum, i) => sum + i.price * i.quantity, 0)
)

onMounted(() => {
  loadAddress()
  loadItems()
})

async function loadAddress() {
  try {
    const list = await getAddressList()
    const def = list.find(a => a.is_default)
    if (def) address.value = def
    else if (list.length > 0) address.value = list[0]
  } catch {}
}

function loadItems() {
  const pages = getCurrentPages()
  const page = pages[pages.length - 1] as any
  const data = page?.$page?.options
  if (data?.items) {
    try {
      items.value = JSON.parse(data.items)
    } catch {}
  }
  if (items.value.length === 0) {
    items.value = [
      { sku_id: 1, product_id: 1, product_name: '示例商品', product_image: '', sku_name: '默认规格', price: 4990, quantity: 1 },
    ]
  }
}

function selectAddress() {
  uni.navigateTo({ url: '/pages/address/address' })
}

function onAddressSelected(addr: Address) {
  address.value = addr
}

function selectCoupon() {
  uni.showToast({ title: '暂无可用优惠券', icon: 'none' })
}

function selectPayment(method: string) {
  payment.value = method
}

async function submitOrder() {
  if (!address.value) {
    uni.showToast({ title: '请选择收货地址', icon: 'none' })
    return
  }
  submitting.value = true
  try {
    const order = await createOrder({
      address_id: address.value.id,
      items: items.value.map(i => ({ sku_id: i.sku_id, quantity: i.quantity })),
      remark: remark.value,
      payment_method: payment.value,
    })
    uni.showToast({ title: '下单成功', icon: 'success' })
    setTimeout(() => {
      uni.navigateTo({ url: `/pages/sub-package-order/orders/orders` })
    }, 1000)
  } catch {
    uni.showToast({ title: '下单失败', icon: 'none' })
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.checkout-page {
  min-height: 100vh;
  background: #f7f8f5;
  display: flex;
  flex-direction: column;
}

.content { flex: 1; }

.section {
  background: #ffffff;
  padding: 28rpx 32rpx;
  border-bottom: 16rpx solid #f7f8f5;
}

.address-card {
  display: flex;
  gap: 20rpx;
  align-items: flex-start;
}

.address-icon {
  width: 72rpx;
  height: 72rpx;
  background: #e8f0e6;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.address-info { flex: 1; }

.address-name-row {
  display: flex;
  align-items: center;
  gap: 16rpx;
  margin-bottom: 8rpx;
}

.address-name {
  font-size: 28rpx;
  font-weight: 600;
  color: #2c2e2c;
}

.address-phone {
  font-size: 24rpx;
  color: #9a9c98;
}

.address-change {
  margin-left: auto;
  font-size: 24rpx;
  color: #9a9c98;
}

.address-detail {
  font-size: 24rpx;
  color: #555755;
  line-height: 1.5;
}

.address-empty {
  font-size: 28rpx;
  color: #bcbeb9;
}

.sec-title {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-bottom: 20rpx;
}

.sec-title-text {
  font-size: 26rpx;
  color: #555755;
  font-weight: 500;
}

.checkout-item {
  display: flex;
  align-items: center;
  gap: 20rpx;
  padding: 12rpx 0;
}

.item-thumb {
  width: 100rpx;
  height: 100rpx;
  border-radius: 12rpx;
  background: #f0f2ed;
  flex-shrink: 0;
}

.item-info {
  flex: 1;
  min-width: 0;
}

.item-name {
  font-size: 26rpx;
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
}

.item-qty {
  font-size: 24rpx;
  color: #9a9c98;
  margin-right: 20rpx;
}

.item-price {
  font-size: 26rpx;
  font-weight: 600;
  color: #d9534f;
}

.coupon-row {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.coupon-label {
  font-size: 28rpx;
  color: #2c2e2c;
}

.coupon-value {
  flex: 1;
  text-align: right;
  font-size: 24rpx;
  color: #9a9c98;
}

.coupon-value.highlight {
  color: #d9534f;
}

.fee-row {
  display: flex;
  justify-content: space-between;
  padding: 12rpx 0;
}

.fee-label {
  font-size: 26rpx;
  color: #555755;
}

.fee-value {
  font-size: 26rpx;
  font-weight: 500;
  color: #2c2e2c;
}

.fee-value.free {
  color: #5b8c5a;
}

.fee-value.discount {
  color: #d9534f;
}

.payment-opt {
  display: flex;
  align-items: center;
  gap: 20rpx;
  padding: 12rpx 0;
}

.radio {
  width: 36rpx;
  height: 36rpx;
  border-radius: 50%;
  border: 4rpx solid #e5e7e3;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.radio.selected {
  border-color: #5b8c5a;
}

.radio-inner {
  width: 20rpx;
  height: 20rpx;
  border-radius: 50%;
  background: #5b8c5a;
}

.pay-name {
  font-size: 28rpx;
  color: #2c2e2c;
}

.remark-row {
  display: flex;
  align-items: center;
  gap: 20rpx;
}

.remark-label {
  font-size: 28rpx;
  color: #2c2e2c;
  flex-shrink: 0;
}

.remark-input {
  flex: 1;
  font-size: 26rpx;
  color: #2c2e2c;
  background: #f0f2ed;
  padding: 16rpx 20rpx;
  border-radius: 12rpx;
}

.bottom-bar {
  height: 120rpx;
  display: flex;
  align-items: center;
  background: #ffffff;
  border-top: 2rpx solid #eef0eb;
  padding: 0 32rpx;
  gap: 20rpx;
  padding-bottom: env(safe-area-inset-bottom, 0);
}

.total-section {
  flex: 1;
}

.total-label {
  font-size: 24rpx;
  color: #555755;
  display: block;
}

.total-amount {
  font-size: 44rpx;
  font-weight: 700;
  color: #d9534f;
}

.total-sym {
  font-size: 28rpx;
}

.btn-submit {
  padding: 22rpx 64rpx;
  background: #5b8c5a;
  color: #fff;
  border-radius: 999rpx;
  font-size: 30rpx;
  font-weight: 600;
  white-space: nowrap;
}

.btn-submit.disabled {
  opacity: 0.5;
}

.placeholder {
  color: #bcbeb9;
  font-size: 26rpx;
}
</style>
