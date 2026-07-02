import { get, post } from './request'

export interface OrderItem {
  product_id: number
  product_name: string
  product_image: string
  sku_id: number
  sku_name: string
  price: number
  quantity: number
}

export interface CreateOrderParams {
  address_id: number
  items: { sku_id: number; quantity: number }[]
  remark?: string
  coupon_id?: number
  payment_method?: string
}

export interface Order {
  id: number
  order_no: string
  status: string
  total_amount: number
  created_at: string
  items: OrderItem[]
}

export const createOrder = (params: CreateOrderParams) =>
  post<Order>('/orders', params, {
    'Idempotent-Key': params.address_id + '-' + Date.now(),
  } as any)

export const getOrderList = (params: { status?: string; page?: number; size?: number }) =>
  get<{ list: Order[]; total: number }>('/orders', params)

export const getOrderDetail = (id: number) =>
  get<Order>(`/orders/${id}`)

export const cancelOrder = (id: number) =>
  post<void>(`/orders/${id}/cancel`)

export const confirmOrder = (id: number) =>
  post<void>(`/orders/${id}/confirm`)
