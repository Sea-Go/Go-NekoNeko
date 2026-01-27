package main

import (
	"context"
	"flag"
	"fmt"
	"sea-try-go/service/points/rpc/internal/mq"
	"time"

	"sea-try-go/service/points/rpc/internal/config"
	"sea-try-go/service/points/rpc/internal/server"
	"sea-try-go/service/points/rpc/internal/svc"
	"sea-try-go/service/points/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/points.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	consumerHandler := mq.NewKafkaConsumer(ctx)
	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		logx.Infof("kafka启动")
		for {
			if err := ctx.KafKa.Consume(rootCtx, []string{c.KafkaConf.Topic}, consumerHandler); err != nil {
				logx.Errorf("Kafka 消费异常: %v", err)
				time.Sleep(time.Second * 3)
			}
			if rootCtx.Err() != nil {
				logx.Infof("kafka 关闭")
				return
			}
		}
	}()
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		__.RegisterPointsServiceServer(grpcServer, server.NewPointsServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
