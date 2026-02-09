// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"sea-try-go/api/internal/config"
	"sea-try-go/api/internal/utils"
	"sea-try-go/service/favorite/favorite_item"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceContext struct {
	Config              config.Config
	Db                  *pgxpool.Pool
	FavoriteItemService *favorite_item.Service
	Cache               *utils.RedisCache
}

func NewServiceContext(c config.Config, db *pgxpool.Pool, cache *utils.RedisCache) *ServiceContext {
	return &ServiceContext{
		Config:              c,
		Db:                  db,
		FavoriteItemService: favorite_item.NewService(db, cache),
		Cache:               cache,
	}
}
