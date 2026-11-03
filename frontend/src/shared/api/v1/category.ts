import { request } from '@/shared/api/http'
import type {
  CategoryBoolResponse,
  CreateCategoryRequest,
  CreateTypeRequest,
  DeleteCategoryRequest,
  DeleteTypeRequest,
  ListCategoriesRequest,
  ListCategoriesResponse,
  ListTypesResponse,
} from '@/shared/types/category'

const auth = (token?: string) => (token ? { token } : {})

export const categoryApi = {
  listTypes(token?: string) {
    return request<ListTypesResponse>('/v1/types/list', { method: 'GET', ...auth(token) })
  },

  listCategories(payload: ListCategoriesRequest, token?: string) {
    return request<ListCategoriesResponse>('/v1/types/category/list', { method: 'POST', body: payload, ...auth(token) })
  },

  createType(payload: CreateTypeRequest, token?: string) {
    return request<CategoryBoolResponse>('/v1/types/create', { method: 'POST', body: payload, ...auth(token) })
  },

  deleteType(payload: DeleteTypeRequest, token?: string) {
    return request<CategoryBoolResponse>('/v1/types/del', { method: 'DELETE', body: payload, ...auth(token) })
  },

  createCategory(payload: CreateCategoryRequest, token?: string) {
    return request<CategoryBoolResponse>('/v1/types/category/create', { method: 'POST', body: payload, ...auth(token) })
  },

  deleteCategory(payload: DeleteCategoryRequest, token?: string) {
    return request<CategoryBoolResponse>('/v1/types/category/del', { method: 'DELETE', body: payload, ...auth(token) })
  },
}
