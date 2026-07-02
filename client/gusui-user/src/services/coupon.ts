import { get } from './request'

export interface Coupon {
  id: number
  name: string
  type: string
  value: number
  min_amount: number
  start_at: string
  end_at: string
  used: boolean
}

export const getCoupons = () =>
  get<Coupon[]>('/coupons')
