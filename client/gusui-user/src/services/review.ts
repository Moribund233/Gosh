import { get } from './request'

export interface Review {
  id: number
  user_name: string
  rating: number
  content: string
  created_at: string
}

export const getReviews = (params: { product_id: number; page?: number; size?: number }) =>
  get<{ list: Review[]; total: number }>('/reviews', params)
