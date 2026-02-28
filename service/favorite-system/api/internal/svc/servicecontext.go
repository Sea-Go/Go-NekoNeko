// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"context"
	"favorite-system/api/internal/config"
	"favorite-system/internal/repo/db"
	"favorite-system/internal/repo/pg"
	"favorite-system/internal/repo/redis"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config config.Config
	DB     *db.Queries
	Redis  *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	pgDB, err := pg.New(context.Background(), c.Postgres.DataSource)
	if err != nil {
		logx.Errorf("failed to connect postgres: %v", err)
		// panic(err) // Optional: panic if DB is critical
	}

	return &ServiceContext{
		Config: c,
		DB:     db.New(pgDB.Pool),
		Redis:  redis.New(c.Redis.Addr),
	}
}
