export interface CategoryType {
  id: number
  name: string
}

export interface Category {
  id: number
  parentID: number
  name: string
}

export interface CreateTypeRequest {
  name: string
}

export interface DeleteTypeRequest {
  id: number
}

export interface ListTypesResponse {
  types: CategoryType[]
}

export interface CreateCategoryRequest {
  parentID: number
  name: string
}

export interface DeleteCategoryRequest {
  id: number
}

export interface ListCategoriesRequest {
  parentID: number
}

export interface ListCategoriesResponse {
  categories: Category[]
}

export interface CategoryBoolResponse {
  success: boolean
}
