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

var configFile = flag.String("f", "etc/favorite.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.Log)

	// 初始化 PostgreSQL 连接池
	pgPool, err := pgxpool.New(context.Background(), c.PgDsn)
	if err != nil {
		logx.Errorf("init pg pool failed: %v", err)
		panic(err)
	}
	defer pgPool.Close()

	// 初始化 Redis 缓存客户端
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
			cache = &utils.RedisCache{} // 使用空的缓存实例（禁用状态）
		} else {
			defer cache.Close()
			logx.Info("redis cache initialized successfully")
		}
	} else {
		cache = &utils.RedisCache{} // 缓存禁用
		logx.Info("redis cache is disabled")
	}

	// 创建 HTTP Server
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 把 pgPool 和 cache 传进 ServiceContext
	ctx := svc.NewServiceContext(c, pgPool, cache)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
