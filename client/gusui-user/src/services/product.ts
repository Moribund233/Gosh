import { get, post } from './request'

export interface Product {
  id: number
  name: string
  price: number
  original_price: number
  images: string[]
  sales: number
  category_id: number
  category_name: string
  is_new: boolean
  is_hot: boolean
  is_featured: boolean
  description: string
  origin: string
  shelf_life: string
  storage_method: string
  net_content: string
  brand: string
  skus: SKU[]
}

export interface SKU {
  id: number
  name: string
  price: number
  stock: number
  specs: Record<string, string>
}

export interface ProductListParams {
  page?: number
  size?: number
  category_id?: number
  keyword?: string
  sort?: string
  is_new?: boolean
  is_hot?: boolean
  is_featured?: boolean
}

export interface ProductListResult {
  list: Product[]
  total: number
}

export const getProducts = (params?: ProductListParams) =>
  get<ProductListResult>('/products', params)

export const getProductDetail = (id: number) =>
  get<Product>(`/products/${id}`)

export const searchProducts = (params: { keyword: string; page?: number; size?: number }) =>
  get<ProductListResult>('/products/search', params)

export const getHotSearch = () =>
  get<string[]>('/products/hot-search')

export const getSearchHistory = () =>
  get<string[]>('/products/search-history')

export const clearSearchHistory = () =>
  post<void>('/products/search-history/clear')
