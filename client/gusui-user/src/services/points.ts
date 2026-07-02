import { get } from './request'

export interface PointsRecord {
  id: number
  amount: number
  type: string
  remark: string
  created_at: string
}

export interface PointsInfo {
  balance: number
  records: PointsRecord[]
}

export const getPoints = () =>
  get<PointsInfo>('/points')
