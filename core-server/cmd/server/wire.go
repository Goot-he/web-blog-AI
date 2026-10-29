//go:build wireinject
// +build wireinject

package main

import (
	"core-server/internal/application"
	"core-server/internal/config"
	"core-server/internal/infras"
	"core-server/internal/rpc"

	"github.com/google/wire"
)

func InitApp(cfg *config.Config) (*rpc.Server, error) {
	wire.Build(
		infras.JobProviderSet,
		application.ProviderSet,
		rpc.ProviderSet,
	)
	return nil, nil
}
