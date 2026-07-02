import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface CartItem {
  id: number
  product_id: number
  product_name: string
  product_image: string
  sku_id: number
  sku_name: string
  price: number
  quantity: number
  stock: number
  selected: boolean
}

export const useCartStore = defineStore('cart', () => {
  const items = ref<CartItem[]>([])

  const selectedItems = computed(() => items.value.filter(i => i.selected))
  const selectedCount = computed(() => selectedItems.value.reduce((sum, i) => sum + i.quantity, 0))
  const selectedAmount = computed(() =>
    selectedItems.value.reduce((sum, i) => sum + i.price * i.quantity, 0)
  )
  const allSelected = computed(() => items.value.length > 0 && items.value.every(i => i.selected))

  function toggleSelect(id: number) {
    const item = items.value.find(i => i.id === id)
    if (item) item.selected = !item.selected
  }

  function toggleSelectAll() {
    const newVal = !allSelected.value
    items.value.forEach(i => { i.selected = newVal })
  }

  function updateQuantity(id: number, quantity: number) {
    const item = items.value.find(i => i.id === id)
    if (item) item.quantity = Math.max(1, Math.min(quantity, item.stock))
  }

  function removeItem(id: number) {
    items.value = items.value.filter(i => i.id !== id)
  }

  function setItems(list: CartItem[]) {
    items.value = list
  }

  function clearCart() {
    items.value = []
  }

  return {
    items, selectedItems, selectedCount, selectedAmount, allSelected,
    toggleSelect, toggleSelectAll, updateQuantity, removeItem, setItems, clearCart
  }
})
