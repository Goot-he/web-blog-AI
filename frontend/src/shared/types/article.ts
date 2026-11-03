export interface Article {
  id: number
  title: string
  summary: string
  content: string
  coverImage: string
  authorID: number
  categoryID: number
  isTop: boolean
  isPublished: boolean
  viewCount: number
  likeCount: number
  favorCount: number
  commentCount: number
  createdAt: number
  updatedAt: number
  publishedAt: number
  authorName: string
  authorAvatar: string
}

export interface ListArticlesRequest {
  page?: number
  pageSize?: number
}

export interface ListByCategoryRequest extends ListArticlesRequest {
  categoryID: number
}

export interface SearchArticlesRequest extends ListArticlesRequest {
  q: string
}

export interface ListArticlesResponse {
  articles: Article[]
  total?: number
}

export interface CreateArticleRequest {
  title: string
  summary?: string
  content: string
  coverImage?: string
  categoryID?: number
  isTop?: boolean
  isPublish?: boolean
  /** 由服务端�?token 注入，无需传�?*/
  authorID?: number
}

export interface EditArticleRequest extends CreateArticleRequest {
  id: number
}

export interface ListMyArticlesRequest extends ListArticlesRequest {
  authorID?: number
  isPublished?: boolean
}

export interface GetArticleRequest {
  id: number
}

export interface GetArticleResponse {
  article: Article
  isLiked: boolean
	isFavored: boolean
}

export interface DeleteArticleRequest {
  id: number
  /** 由服务端�?token 注入，无需传�?*/
  authorID?: number
}

export interface PublishDraftRequest {
  id: number
}

export interface DeleteDraftRequest {
  id: number
}

export interface ArticleBoolResponse {
  success: boolean
}

export interface EditArticleResponse {
  success: boolean
  articleID: number
}

export interface ArticleImageUploadResponse {
  url: string
}
