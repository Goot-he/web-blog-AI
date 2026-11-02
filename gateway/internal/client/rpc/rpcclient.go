package rpc

import (
	"fmt"
	"gateway/internal/client/rpc/core-rpc/articlepb"
	"gateway/internal/client/rpc/core-rpc/authpb"
	"gateway/internal/client/rpc/core-rpc/categorypb"
	"gateway/internal/client/rpc/core-rpc/commentpb"
	"gateway/internal/client/rpc/core-rpc/favorpb"
	"gateway/internal/client/rpc/core-rpc/likepb"
	"gateway/internal/client/rpc/core-rpc/userpb"
	"gateway/internal/config"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	CoreConn *grpc.ClientConn

	AuthClient     authpb.AuthServiceClient
	UserClient     userpb.UserServiceClient
	ArticleClient  articlepb.ArticleServiceClient
	CommentClient  commentpb.CommentServiceClient
	CategoryClient categorypb.CategoryServiceClient
	LikeClient     likepb.LikeServiceClient
	FavorClient    favorpb.FavorServiceClient

	// 请求超时时间
	requestTimeout time.Duration
}

func NewRpcClient(cfg *config.Config) (*Client, error) {
	coreConn, err := newCoreConn(cfg)
	if err != nil {
		return nil, err
	}

	timeoutSec := cfg.Rpc.RequestTimeout
	if timeoutSec == 0 {
		timeoutSec = 5
	}

	return &Client{
		CoreConn:       coreConn,
		AuthClient:     authpb.NewAuthServiceClient(coreConn),
		UserClient:     userpb.NewUserServiceClient(coreConn),
		ArticleClient:  articlepb.NewArticleServiceClient(coreConn),
		CommentClient:  commentpb.NewCommentServiceClient(coreConn),
		CategoryClient: categorypb.NewCategoryServiceClient(coreConn),
		LikeClient:     likepb.NewLikeServiceClient(coreConn),
		FavorClient:    favorpb.NewFavorServiceClient(coreConn),
		requestTimeout: time.Duration(timeoutSec) * time.Second,
	}, nil
}

// 创建一个到core-server grpc服务的客户端连接对象，交给生成的客户端使用
func newCoreConn(cfg *config.Config) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		cfg.Rpc.CoreServerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("dial core-server err %s: %v", cfg.Rpc.CoreServerAddr, err)
	}

	return conn, nil
}
