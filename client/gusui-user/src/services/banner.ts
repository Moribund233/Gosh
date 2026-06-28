import { get } from './request'

export interface Banner {
  id: number
  title: string
  image: string
  link: string
  sort_order: number
}

export interface BrandStory {
  id: number
  title: string
  content: string
  image: string
}

export const getBanners = () =>
  get<Banner[]>('/banners')

export const getBrandStory = () =>
  get<BrandStory>('/brand-story')
