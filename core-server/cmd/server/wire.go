//go:build wireinject
// +build wireinject

package main

import (
	"core-server/internal/application"
	"core-server/internal/config"
	"core-server/internal/infras/repo"
	"core-server/internal/rpc"

	"github.com/google/wire"
)

func InitApp(cfg *config.Config) (*App, error) {
	wire.Build(
		repo.SetRepo,
		application.SetService,
		rpc.SetRpc,
		wire.Struct(new(App), "*"),
	)
	return nil, nil
}
