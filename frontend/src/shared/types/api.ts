export interface ApiResponse<T> {
  code: number
  msg: string
  data: T
}

export interface PageRequest {
  page: number
  pageSize: number
}

export interface PageResponse<T> {
  list: T[]
  total: number
  page: number
  pageSize: number
}
