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

	"sea-try-go/service/common/logger"
	"sea-try-go/service/user/user/api/internal/config"
	"sea-try-go/service/user/user/api/internal/handler"
	"sea-try-go/service/user/user/api/internal/svc"
)

var configFile = flag.String("f", "etc/usercenter.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.Log)
	logger.Init(c.Name)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	pgPool, err := pgxpool.New(context.Background(), c.PgDsn)
	if err != nil {
		logx.Errorf("init pg pool failed: %v", err)
		panic(err)
	}
	defer pgPool.Close()

	ctx := svc.NewServiceContext(c, pgPool)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)

	server.Start()
}
