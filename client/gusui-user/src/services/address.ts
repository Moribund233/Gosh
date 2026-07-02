import { get, post, put, del } from './request'

export interface Address {
  id: number
  name: string
  phone: string
  province: string
  city: string
  district: string
  detail: string
  is_default: boolean
}

export const getAddressList = () =>
  get<Address[]>('/addresses')

export const getAddressDetail = (id: number) =>
  get<Address>(`/addresses/${id}`)

export const createAddress = (data: Omit<Address, 'id'>) =>
  post<Address>('/addresses', data)

export const updateAddress = (id: number, data: Partial<Address>) =>
  put<Address>(`/addresses/${id}`, data)

export const deleteAddress = (id: number) =>
  del(`/addresses/${id}`)

export const setDefaultAddress = (id: number) =>
  put<void>(`/addresses/${id}/default`)
