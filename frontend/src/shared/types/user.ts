import type { AuthUser } from '@/shared/types/auth'

export interface CurrentUser extends AuthUser {
  email: string
  phone: string
  sex: string
  birthday: string
  signature: string
  article_count: number
  view_count: number
  comment_count: number
  followers_count: number
  following_count: number
  like_count: number
  receive_like_count: number
  favor_count: number
  receive_favor_count: number
}

export interface VisitorProfile {
  id: number
  username: string
  avatar: string
  sex: string
  birthday: string
  signature: string
  article_count: number
  view_count: number
  comment_count: number
  followers_count: number
  following_count: number
  like_count: number
  receive_like_count: number
  favor_count: number
  receive_favor_count: number
}

export interface ManagedUser {
  id: number
  username: string
  email: string
  phone: string
  avatar: string
  sex: string
  birthday: string
  signature: string
  role: string
  status: 'approved' | 'blocked'
}

export interface UserListResponse {
  users: ManagedUser[]
  total: number
}

export interface UserAvatarResponse {
  avatar: string
}
