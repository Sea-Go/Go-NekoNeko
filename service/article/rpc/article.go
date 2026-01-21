package main

import (
	"flag"
	"fmt"
	"sea-try-go/service/article/rpc/internal/config"
	"sea-try-go/service/article/rpc/internal/model"
	"sea-try-go/service/article/rpc/internal/server"
	"sea-try-go/service/article/rpc/internal/svc"
	"sea-try-go/service/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/article.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	u := model.NewArticleRepo(c)
	ctx := svc.NewServiceContext(c, u)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		__.RegisterArticleServiceServer(grpcServer, server.NewArticleServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
