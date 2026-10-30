package rpc

import (
	"context"
	"core-server/internal/application"
	"core-server/internal/model/dto"
	"core-server/internal/rpc/articlepb"
	"core-server/internal/rpc/likepb"
)

type LikeRpc struct {
	likepb.UnimplementedLikeServiceServer
	service *application.LikeService
}

func NewLikeRpc(likeService *application.LikeService) *LikeRpc {
	return &LikeRpc{
		service: likeService,
	}
}

func (r *LikeRpc) HasLike(ctx context.Context, req *likepb.HasLikeRequest) (*likepb.HasLikeResponse, error) {
	isThumbUp, err := r.service.HasThumbUp(ctx, req.GetUserID(), req.GetObjectID(), req.GetObjectType())
	if err != nil {
		return nil, err
	}
	return &likepb.HasLikeResponse{
		IsLiked: isThumbUp,
	}, nil
}

func (r *LikeRpc) ThumbUp(ctx context.Context, req *likepb.ThumbUpRequest) (*likepb.ThumbUpResponse, error) {
	err := r.service.ThumbUp(ctx, req.GetUserID(), req.GetObjectID(), req.GetObjectType())
	if err != nil {
		return nil, err
	}
	return &likepb.ThumbUpResponse{
		Success: true,
	}, nil
}

func (r *LikeRpc) CancelThumbUp(ctx context.Context, req *likepb.CancelThumbUpRequest) (*likepb.CancelThumbUpResponse, error) {
	err := r.service.CancelThumbUp(ctx, req.GetUserID(), req.GetObjectID(), req.GetObjectType())
	if err != nil {
		return nil, err
	}
	return &likepb.CancelThumbUpResponse{
		Success: true,
	}, nil
}

func (r *LikeRpc) PageQueryUserLikeList(ctx context.Context, req *likepb.PageQueryUserLikeListRequest) (*likepb.PageQueryUserLikeListResponse, error) {
	articles, err := r.service.UserLikeList(ctx, req.GetUserID(), req.GetObjectType(), int(req.Page), int(req.PageSize))

	if err != nil {
		return nil, err
	}

	items := make([]*articlepb.Article, 0, len(articles))
	for _, article := range articles {
		if article == nil || article.Article == nil {
			continue
		}
		items = append(items, toArticlePB(dto.ToArticleResponse(article)))
	}

	return &likepb.PageQueryUserLikeListResponse{
		Articles: items,
		Total:    int64(len(items)),
	}, nil
}
