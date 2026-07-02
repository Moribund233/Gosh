<template>
  <view class="edit-page">
    <scroll-view class="content" scroll-y>
      <view class="form-group">
        <text class="form-label">收货人</text>
        <input class="form-input" v-model="form.name" type="text" placeholder="请输入收货人姓名" placeholder-class="placeholder" />
      </view>

      <view class="form-group">
        <text class="form-label">手机号码</text>
        <input class="form-input" v-model="form.phone" type="text" maxlength="11" placeholder="请输入手机号码" placeholder-class="placeholder" />
      </view>

      <view class="form-group">
        <text class="form-label">所在地区</text>
        <view class="region-picker" @tap="showRegionPicker">
          <text class="region-text" :class="{ muted: !form.province }">{{ regionText || '请选择省/市/区' }}</text>
          <tui-icon name="arrowright" size="28" color="#bcbeb9"></tui-icon>
        </view>
      </view>

      <view class="form-group">
        <text class="form-label">详细地址</text>
        <input class="form-input" v-model="form.detail" type="text" placeholder="街道、门牌号等" placeholder-class="placeholder" />
      </view>

      <view class="form-group switch-group">
        <text class="form-label">设为默认地址</text>
        <switch class="form-switch" :checked="form.is_default" color="#5b8c5a" @change="onDefaultChange" />
      </view>

      <view class="submit-wrapper">
        <view class="btn-submit" :disabled="!canSubmit" @tap="handleSubmit">
          <text class="btn-text">保存</text>
        </view>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getAddressDetail, createAddress, updateAddress, type Address } from '@/services/address'

const addressId = ref<number | null>(null)
const form = ref({
  name: '',
  phone: '',
  province: '',
  city: '',
  district: '',
  detail: '',
  is_default: false,
})

const regionText = computed(() => {
  const parts = [form.value.province, form.value.city, form.value.district]
  return parts.filter(Boolean).join(' ')
})

const canSubmit = computed(() =>
  form.value.name && form.value.phone.length === 11 && form.value.province && form.value.detail
)

onMounted(() => {
  const pages = getCurrentPages()
  const page = pages[pages.length - 1] as any
  const id = page?.$page?.options?.id || page?.options?.id
  if (id) {
    addressId.value = Number(id)
    loadAddress(Number(id))
  }
})

async function loadAddress(id: number) {
  try {
    const data = await getAddressDetail(id)
    form.value = {
      name: data.name,
      phone: data.phone,
      province: data.province,
      city: data.city,
      district: data.district,
      detail: data.detail,
      is_default: data.is_default,
    }
  } catch {}
}

function showRegionPicker() {
  uni.showActionSheet({
    itemList: ['请通过城市选择器选择', '（实际项目中使用 picker 组件）'],
    success: () => {},
  })
}

function onDefaultChange(e: any) {
  form.value.is_default = e.detail.value
}

async function handleSubmit() {
  if (!canSubmit.value) return
  try {
    if (addressId.value) {
      await updateAddress(addressId.value, form.value)
    } else {
      await createAddress(form.value as any)
    }
    uni.showToast({ title: '保存成功', icon: 'success' })
    setTimeout(() => uni.navigateBack(), 1000)
  } catch {
    uni.showToast({ title: '保存失败', icon: 'none' })
  }
}
</script>

<style scoped>
.edit-page {
  min-height: 100vh;
  background: #f7f8f5;
}

.content {
  height: 100vh;
  padding: 20rpx 0;
}

.form-group {
  padding: 24rpx 32rpx;
  background: #ffffff;
  margin-bottom: 2rpx;
}

.form-label {
  font-size: 26rpx;
  color: #555755;
  font-weight: 500;
  display: block;
  margin-bottom: 16rpx;
}

.form-input {
  width: 100%;
  height: 72rpx;
  padding: 0 16rpx;
  background: #f0f2ed;
  border-radius: 12rpx;
  font-size: 28rpx;
  color: #2c2e2c;
  border: 2rpx solid transparent;
  box-sizing: border-box;
}

.form-input:focus {
  border-color: #5b8c5a;
  background: #ffffff;
}

.region-picker {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 72rpx;
  padding: 0 16rpx;
  background: #f0f2ed;
  border-radius: 12rpx;
}

.region-text {
  font-size: 28rpx;
  color: #2c2e2c;
}

.region-text.muted {
  color: #bcbeb9;
}

.switch-group {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.switch-group .form-label {
  margin-bottom: 0;
}

.form-switch {
  transform: scale(0.8);
}

.submit-wrapper {
  padding: 40rpx 32rpx 60rpx;
}

.btn-submit {
  width: 100%;
  height: 88rpx;
  background: #5b8c5a;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-submit:active {
  opacity: 0.85;
}

.btn-text {
  font-size: 30rpx;
  font-weight: 600;
  color: #ffffff;
}

.placeholder {
  color: #bcbeb9;
  font-size: 28rpx;
}
</style>
