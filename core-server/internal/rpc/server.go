package rpc

import (
	"core-server/internal/infras/repo"
	"core-server/internal/rpc/articlepb"
	"core-server/internal/rpc/authpb"
	"core-server/internal/rpc/userpb"
	"fmt"
	"log"
	"net"

	"core-server/internal/config"
	"core-server/internal/infras/cache"
	"core-server/internal/rpc/categorypb"
	"core-server/internal/rpc/commentpb"
	"core-server/internal/rpc/likepb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	cfg         *config.Config
	Engine      *grpc.Server
	dbClient    *repo.EntClient
	cacheClient *cache.CacheClient
}

func NewServer(
	cfg *config.Config,
	dbClient *repo.EntClient,
	cacheClient *cache.CacheClient,
	auth *AuthRPC,
	like *LikeRPC,
	category *CategoryRpc,
	article *ArticleRpc,
	comment *CommentRpc,
	user *UserRpc,
) (*Server, error) {
	engine := grpc.NewServer()
	authpb.RegisterAuthServiceServer(engine, auth)
	likepb.RegisterLikeServiceServer(engine, like)
	categorypb.RegisterCategoryServiceServer(engine, category)
	articlepb.RegisterArticleServiceServer(engine, article)
	commentpb.RegisterCommentServiceServer(engine, comment)
	userpb.RegisterUserServiceServer(engine, user)
	reflection.Register(engine)

	return &Server{
		cfg:         cfg,
		Engine:      engine,
		dbClient:    dbClient,
		cacheClient: cacheClient,
	}, nil
}

func (s *Server) Start() error {

	lis, err := net.Listen("tcp", s.cfg.ServerConf.Address)
	if err != nil {
		return fmt.Errorf("listen %s: %w", s.cfg.ServerConf.Address, err)
	}
	log.Printf("core-server listening on %s", s.cfg.ServerConf.Address)
	return s.Engine.Serve(lis)
}

func (s *Server) Stop() {
	s.Engine.GracefulStop()

	if s.dbClient != nil {
		if err := s.dbClient.Close(); err != nil {
			log.Printf("close mysql: %v", err)
		}
	}
	if s.cacheClient != nil {
		if err := s.cacheClient.Close(); err != nil {
			log.Printf("close redis: %v", err)
		}
	}
}

func (s *Server) Addr() string {
	return s.cfg.ServerConf.Address
}
