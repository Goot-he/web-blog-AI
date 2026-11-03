import { request, upload } from '@/shared/api/http'
import type {
  ArticleImageUploadResponse,
  ArticleBoolResponse,
  CreateArticleRequest,
  DeleteArticleRequest,
  DeleteDraftRequest,
  EditArticleRequest,
  EditArticleResponse,
  GetArticleRequest,
  GetArticleResponse,
  ListArticlesRequest,
  ListArticlesResponse,
  ListByCategoryRequest,
  ListMyArticlesRequest,
  PublishDraftRequest,
  SearchArticlesRequest,
} from '@/shared/types/article'

const auth = (token?: string) => (token ? { token } : {})

export const articleApi = {
  list(params: ListArticlesRequest = {}, token?: string) {
    return request<ListArticlesResponse>('/v1/article/list', { method: 'POST', body: params, ...auth(token) })
  },

  byCategory(categoryID: number, params: ListArticlesRequest = {}, token?: string) {
    const body: ListByCategoryRequest = { categoryID, ...params }
    return request<ListArticlesResponse>('/v1/article/list/by_cate', { method: 'POST', body, ...auth(token) })
  },

  search(q: string, params: ListArticlesRequest = {}, token?: string) {
    const body: SearchArticlesRequest = { q, ...params }
    return request<ListArticlesResponse>('/v1/article/search', { method: 'POST', body, ...auth(token) })
  },

  detail(payload: GetArticleRequest, token?: string) {
    return request<GetArticleResponse>('/v1/article/message', { method: 'POST', body: payload, ...auth(token) })
  },

  create(payload: CreateArticleRequest, token?: string) {
    return request<ArticleBoolResponse>('/v1/article/create', { method: 'POST', body: payload, ...auth(token) })
  },

  saveDraft(payload: CreateArticleRequest, token?: string) {
    return request<ArticleBoolResponse>('/v1/article/draft', { method: 'POST', body: payload, ...auth(token) })
  },

  edit(payload: EditArticleRequest, token?: string) {
    return request<EditArticleResponse>('/v1/article/edit', { method: 'POST', body: payload, ...auth(token) })
  },

  listByUser(params: ListMyArticlesRequest = {}, token?: string | null) {
    return request<ListArticlesResponse>('/v1/article/list/by_user', { method: 'POST', body: params, ...(token === null ? { auth: false } : auth(token)) })
  },

  listDrafts(params: ListArticlesRequest = {}, token?: string) {
    return request<ListArticlesResponse>('/v1/article/draft/list', { method: 'POST', body: params, ...auth(token) })
  },

  publishDraft(payload: PublishDraftRequest, token?: string) {
    return request<ArticleBoolResponse>('/v1/article/draft/publish', { method: 'POST', body: payload, ...auth(token) })
  },

  uploadCover(file: File) {
    return upload<ArticleImageUploadResponse>('/v1/article/upload/cover', file)
  },

  uploadContent(file: File) {
    return upload<ArticleImageUploadResponse>('/v1/article/upload/content', file)
  },

  delete(payload: DeleteArticleRequest, token?: string) {
    return request<ArticleBoolResponse>('/v1/article/del', { method: 'DELETE', body: payload, ...auth(token) })
  },

  deleteDraft(payload: DeleteDraftRequest, token?: string) {
    return request<ArticleBoolResponse>('/v1/article/draft/del', { method: 'DELETE', body: payload, ...auth(token) })
  },
}
