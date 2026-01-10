package main

import (
	"flag"
	"fmt"
	"sea-try-go/service/favorite/rpc/internal/model/postgres"

	"sea-try-go/service/favorite/rpc/internal/config"
	"sea-try-go/service/favorite/rpc/internal/server"
	"sea-try-go/service/favorite/rpc/internal/svc"
	"sea-try-go/service/favorite/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/favorite.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	u := postgres.NewFavoriteService(c)
	ctx := svc.NewServiceContext(c, u)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		__.RegisterFavoriteServiceServer(grpcServer, server.NewFavoriteServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
