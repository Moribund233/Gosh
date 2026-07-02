import { post, get, put } from './request'

export interface LoginParams {
  phone: string
  password: string
}

export interface RegisterParams {
  phone: string
  password: string
  code: string
}

export interface AuthResult {
  token: string
  user: UserInfo
}

export interface UserInfo {
  id: number
  phone: string
  nickname: string
  avatar: string
  created_at: string
}

export const login = (params: LoginParams) =>
  post<AuthResult>('/auth/login', params)

export const register = (params: RegisterParams) =>
  post<AuthResult>('/auth/register', params)

export const sendCode = (phone: string) =>
  post<void>('/auth/send-code', { phone })

export const getUserInfo = () =>
  get<UserInfo>('/user/info')

export const updateUserInfo = (data: Partial<UserInfo>) =>
  put<UserInfo>('/user/info', data)
