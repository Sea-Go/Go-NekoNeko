package svc

import (
	"sea-try-go/service/favorite/rpc/internal/config"
	"sea-try-go/service/favorite/rpc/internal/model/postgres"
)

type ServiceContext struct {
	Config       config.Config
	FavoriteRepo *postgres.FavoriteRepo
}

func NewServiceContext(c config.Config, favoriteRepo *postgres.FavoriteRepo) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		FavoriteRepo: favoriteRepo,
	}
}
