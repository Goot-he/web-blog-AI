import { request } from '@/shared/api/http'
import type { LikeBoolResponse, LikeRequest, UserLikeListRequest, UserLikeListResponse } from '@/shared/types/like'

const auth = (token?: string) => (token ? { token } : {})

export const likeApi = {
  thumbUp(payload: LikeRequest, token?: string) {
    return request<LikeBoolResponse>('/v1/like/thumb_up', { method: 'POST', body: payload, ...auth(token) })
  },

  cancelThumbUp(payload: LikeRequest, token?: string) {
    return request<LikeBoolResponse>('/v1/like/cancel_thumb_up', { method: 'POST', body: payload, ...auth(token) })
  },

  list(payload: UserLikeListRequest, token?: string) {
    return request<UserLikeListResponse>('/v1/like/list', { method: 'POST', body: payload, ...auth(token) })
  },
}
