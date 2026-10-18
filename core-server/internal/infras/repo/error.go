package repo

import "errors"

var (
	ErrorNotFound = errors.New("resources not found")

	ErrorArticleNotFound  = errors.New("article not found")
	ErrorUserNotFound     = errors.New("user not found")
	ErrorCategoryNotFound = errors.New("category not found")
	ErrorDraftNotFound    = errors.New("draft not found")

	// ErrorNotDraft 已发布的文章不能走草稿相关的操作
	ErrorNotDraft           = errors.New("article is not a draft")
	ErrorArticleIsPublished = errors.New("article is published")

	ErrorCommentNotFound = errors.New("comment not found")
	ErrorNoPermission    = errors.New("permission denied")
	ErrorInvalidParam    = errors.New("invalid parameter")
)
