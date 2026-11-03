import { request } from '@/shared/api/http'
import type { FavorBoolResponse, FavorRequest, UserFavorListRequest, UserFavorListResponse, FavoriteFolder, FavoriteFolderListResponse } from '@/shared/types/favor'

const auth = (token?: string) => token ? { token } : {}

export const favorApi = {
  favor(payload: FavorRequest, token?: string) { return request<FavorBoolResponse>('/v1/favor', { method: 'POST', body: payload, ...auth(token) }) },
  cancel(payload: FavorRequest, token?: string) { return request<FavorBoolResponse>('/v1/favor/cancel', { method: 'POST', body: payload, ...auth(token) }) },
  list(payload: UserFavorListRequest, token?: string) { return request<UserFavorListResponse>('/v1/favor/list', { method: 'POST', body: payload, ...auth(token) }) },
  folders(token?: string) { return request<FavoriteFolderListResponse>('/v1/favor/folders', { ...auth(token) }) },
  createFolder(name: string, token?: string) { return request<FavoriteFolder>('/v1/favor/folders', { method: 'POST', body: { name }, ...auth(token) }) },
  deleteFolder(folderId: number, token?: string) { return request<FavorBoolResponse>('/v1/favor/folders/delete', { method: 'POST', body: { folderId }, ...auth(token) }) },
}
