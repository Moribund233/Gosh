import { get } from './request'

export interface FlashSale {
  id: number
  product_id: number
  product_name: string
  product_image: string
  flash_price: number
  original_price: number
  flash_stock: number
  start_at: string
  end_at: string
}

export const getFlashSales = () =>
  get<FlashSale[]>('/flash-sales')
