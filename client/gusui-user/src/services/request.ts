const BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

interface RequestOptions {
  url: string
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  data?: any
  header?: Record<string, string>
}

interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export function request<T = any>(options: RequestOptions): Promise<T> {
  return new Promise((resolve, reject) => {
    const token = uni.getStorageSync('token')
    const header: Record<string, string> = {
      'Content-Type': 'application/json',
      ...options.header,
    }
    if (token) {
      header['Authorization'] = `Bearer ${token}`
    }

    uni.request({
      url: BASE_URL + options.url,
      method: options.method || 'GET',
      data: options.data,
      header,
      success: (res) => {
        const data = res.data as ApiResponse<T>
        if (res.statusCode === 200) {
          resolve(data.data)
        } else if (res.statusCode === 401) {
          uni.removeStorageSync('token')
          uni.removeStorageSync('userInfo')
          uni.navigateTo({ url: '/pages/login/login' })
          reject(new Error('未登录'))
        } else {
          uni.showToast({ title: data.message || '请求失败', icon: 'none' })
          reject(new Error(data.message))
        }
      },
      fail: (err) => {
        uni.showToast({ title: '网络开小差', icon: 'none' })
        reject(err)
      },
    })
  })
}

export const get = <T>(url: string, data?: any) =>
  request<T>({ url, method: 'GET', data })

export const post = <T>(url: string, data?: any) =>
  request<T>({ url, method: 'POST', data })

export const put = <T>(url: string, data?: any) =>
  request<T>({ url, method: 'PUT', data })

export const del = <T>(url: string) =>
  request<T>({ url, method: 'DELETE' })
