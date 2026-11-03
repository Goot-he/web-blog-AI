import { request } from '@/shared/api/http'
import type {
  CommentBoolResponse,
  CommentListResponse,
  CommentRepliesResponse,
  CreateCommentRequest,
  CreateReplyRequest,
  DeleteCommentRequest,
  DeleteReplyRequest,
  GetArticleCommentsRequest,
  GetCommentRepliesRequest,
} from '@/shared/types/comment'

const auth = (token?: string) => (token ? { token } : {})

export const commentApi = {
  list(payload: GetArticleCommentsRequest, token?: string) {
    return request<CommentListResponse>('/v1/comment/list', { method: 'POST', body: payload, ...auth(token) })
  },

  replies(payload: GetCommentRepliesRequest, token?: string) {
    return request<CommentRepliesResponse>('/v1/comment/replies', { method: 'POST', body: payload, ...auth(token) })
  },

  create(payload: CreateCommentRequest, token?: string) {
    return request<CommentBoolResponse>('/v1/comment/create', { method: 'POST', body: payload, ...auth(token) })
  },

  reply(payload: CreateReplyRequest, token?: string) {
    return request<CommentBoolResponse>('/v1/comment/reply', { method: 'POST', body: payload, ...auth(token) })
  },

  delete(payload: DeleteCommentRequest, token?: string) {
    return request<CommentBoolResponse>('/v1/comment/delete', { method: 'DELETE', body: payload, ...auth(token) })
  },

  deleteReply(payload: DeleteReplyRequest, token?: string) {
    return request<CommentBoolResponse>('/v1/comment/reply/delete', { method: 'DELETE', body: payload, ...auth(token) })
  },
}
