import { get, post, put, del } from './request'

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

export const getCartList = () =>
  get<CartItem[]>('/cart')

export const addToCart = (data: { product_id: number; sku_id: number; quantity: number }) =>
  post<CartItem>('/cart', data)

export const updateCartQuantity = (id: number, quantity: number) =>
  put<void>(`/cart/${id}`, { quantity })

export const toggleCartItem = (id: number, selected: boolean) =>
  put<void>(`/cart/${id}/select`, { selected })

export const toggleCartAll = (selected: boolean) =>
  put<void>('/cart/select-all', { selected })

export const removeCartItem = (id: number) =>
  del(`/cart/${id}`)

export const getCartCount = () =>
  get<{ count: number }>('/cart/count')
