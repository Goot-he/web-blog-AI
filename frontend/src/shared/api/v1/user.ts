import { request, upload } from '@/shared/api/http'
import type { CurrentUser, UserAvatarResponse, UserListResponse, VisitorProfile } from '@/shared/types/user'

export interface UpdateProfilePayload {
  username: string
  phone: string
  sex: '' | 'male' | 'female'
  birthday: string
  signature: string
}

export const userApi = {
  current(token?: string) {
    return request<CurrentUser>('/v1/user/profile', { method: 'GET', ...(token ? { token } : {}) })
  },

  getProfile() {
    return request<CurrentUser>('/v1/user/profile')
  },

  getPublicProfile(id: number) {
    return request<VisitorProfile>(`/v1/user/profile/${id}`, { auth: false })
  },

  updateProfile(payload: UpdateProfilePayload) {
    return request<null>('/v1/user/profile', { method: 'PUT', body: payload }).then(() => userApi.getProfile())
  },

  updateAvatar(file: File) {
    return upload<UserAvatarResponse>('/v1/user/avatar', file, {}, { method: 'PUT' })
  },

  list(keyword = '', page = 1, pageSize = 20) {
    const search = keyword.trim()
    if (search) {
      const query = new URLSearchParams({ keyword: search, page: String(page), page_size: String(pageSize) })
      return request<UserListResponse>(`/v1/admin/search?${query.toString()}`)
    }
    const query = new URLSearchParams({ page: String(page), pageSize: String(pageSize) })
    return request<UserListResponse>(`/v1/admin/list?${query.toString()}`)
  },

  updateBlacklist(userID: number, blacklisted: boolean) {
    const path = blacklisted ? '/v1/admin/blacklist/add' : '/v1/admin/blacklist/remove'
    return request<null>(path, { method: 'PUT', body: { user_id: userID } })
  },
}
