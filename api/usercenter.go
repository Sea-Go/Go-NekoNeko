//go:build tools
// +build tools

// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"

	"sea-try-go/api/internal/config"
	"sea-try-go/api/internal/handler"
	"sea-try-go/api/internal/svc"
	"sea-try-go/api/internal/utils"
)

var configFile = flag.String("f", "etc/usercenter.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.Log)

	// Initialize PostgreSQL connection pool
	pgPool, err := pgxpool.New(context.Background(), c.PgDsn)
	if err != nil {
		logx.Errorf("init pg pool failed: %v", err)
		panic(err)
	}
	defer pgPool.Close()

	// Initialize Redis cache client
	var cache *utils.RedisCache
	if c.Redis.Enabled {
		cache, err = utils.NewRedisCache(utils.CacheConfig{
			Addr:     c.Redis.Addr,
			Password: c.Redis.Password,
			TTL:      c.Redis.TTL,
			Enabled:  c.Redis.Enabled,
		})
		if err != nil {
			logx.Error(fmt.Sprintf("init redis cache failed: %v, will run without cache", err))
			cache = &utils.RedisCache{} // Use empty cache instance (disabled)
		} else {
			defer cache.Close()
			logx.Info("redis cache initialized successfully")
		}
	} else {
		cache = &utils.RedisCache{} // Cache disabled
		logx.Info("redis cache is disabled")
	}

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c, pgPool, cache)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
