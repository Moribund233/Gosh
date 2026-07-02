<script setup lang="ts">
import { onLaunch } from '@dcloudio/uni-app'

const authPages = [
  '/pages/address/address',
  '/pages/address-edit/address-edit',
  '/pages/coupon/coupon',
  '/pages/points/points',
  '/pages/cart/cart',
]

const subAuthPages = [
  '/pages/sub-package-order/checkout/checkout',
  '/pages/sub-package-order/orders/orders',
]

onLaunch(() => {
  const token = uni.getStorageSync('token')
  if (token) {
    console.log('App Launch - 已登录')
  }

  uni.addInterceptor('navigateTo', { invoke: checkAuth })
  uni.addInterceptor('redirectTo', { invoke: checkAuth })
  uni.addInterceptor('switchTab', { invoke: checkAuth })
  uni.addInterceptor('reLaunch', { invoke: checkAuth })
})

function checkAuth(args: { url: string }) {
  const token = uni.getStorageSync('token')
  const url = args.url.split('?')[0]

  const needsAuth = authPages.some(p => url.startsWith(p)) ||
    subAuthPages.some(p => url.startsWith(p))

  if (needsAuth && !token) {
    uni.navigateTo({ url: '/pages/login/login' })
    return false
  }
}
</script>

<style>
page {
  background-color: #f7f8f5;
  font-family: -apple-system, BlinkMacSystemFont, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;
  font-size: 14px;
  color: #2c2e2c;
  line-height: 1.5;
}

uni-tabbar { display: none !important; }

.flex { display: flex; }
.flex-col { flex-direction: column; }
.items-center { align-items: center; }
.justify-center { justify-content: center; }
.justify-between { justify-content: space-between; }
.gap-1 { gap: 4px; }
.gap-2 { gap: 8px; }
.gap-3 { gap: 12px; }
.gap-4 { gap: 16px; }
.text-center { text-align: center; }
.truncate { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
</style>
