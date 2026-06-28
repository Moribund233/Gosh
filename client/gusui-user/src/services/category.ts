import { get } from './request'

export interface Category {
  id: number
  name: string
  icon: string
  parent_id: number | null
  children: Category[]
}

export const getCategoryTree = () =>
  get<Category[]>('/categories')

export const getCategoryById = (id: number) =>
  get<Category>(`/categories/${id}`)
