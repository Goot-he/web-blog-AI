import type { Article } from '@/shared/types/article'

export interface LikeRequest {
  objectType: string
  objectId: number
  /** 由服务端�?token 注入，无需传�?*/
  userId?: number
}

export interface UserLikeListRequest {
  objectType: string
  page?: number
  pageSize?: number
  /** 由服务端�?token 注入，无需传�?*/
  userId?: number
}

export interface LikeBoolResponse {
  success: boolean
}

export interface UserLikeListResponse {
  articles: Article[]
  total: number
}
