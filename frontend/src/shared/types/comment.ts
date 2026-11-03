export interface CommentItem {
  id: number
  articleId: number
  userId: number
  parentId: number
  rootId: number
  replyToId: number
  replyToUserName: string
  content: string
  likeCount: number
  childCount: number
  createdAt: string
  userName: string
  userAvatar: string
  isLiked: boolean
}

export interface CreateCommentRequest {
  articleId: number
  content: string
  /** 由服务端�?token 注入，无需传�?*/
  userId?: number
}

export interface CreateReplyRequest {
  articleId: number
  parentId: number
  replyToId?: number
  content: string
  /** 由服务端�?token 注入，无需传�?*/
  userId?: number
}

export interface DeleteCommentRequest {
  id: number
  articleId: number
  /** 由服务端�?token 注入，无需传�?*/
  userId?: number
}

export interface DeleteReplyRequest {
  id: number
  articleId: number
  /** 由服务端从 token 注入，无需传入 */
  userId?: number
}

export interface GetArticleCommentsRequest {
  articleId: number
  page?: number
  size?: number
}

export interface GetCommentRepliesRequest {
  parentId: number
  page?: number
  size?: number
}

export interface CommentListResponse {
  comments: CommentItem[]
  page: number
  size: number
}

export interface CommentRepliesResponse {
  replies: CommentItem[]
  page: number
  size: number
}

export interface CommentBoolResponse {
  success: boolean
}
