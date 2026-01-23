package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/zeromicro/go-queue/kq"
	"os"
	"sea-try-go/service/article/rpc/internal/config"
	"sea-try-go/service/article/rpc/internal/model"
	"sea-try-go/service/article/rpc/internal/mqs"
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

	if c.AliGreen.AccessKeyId == "" {
		c.AliGreen.AccessKeyId = os.Getenv("ALIYUN_ACCESS_KEY_ID")
	}
	if c.AliGreen.AccessKeySecret == "" {
		c.AliGreen.AccessKeySecret = os.Getenv("ALIYUN_ACCESS_KEY_SECRET")
	}

	// 3. 校验密钥是否配置（可选，防止程序启动后报错）
	if c.AliGreen.AccessKeyId == "" || c.AliGreen.AccessKeySecret == "" {
		panic("环境变量 ALIYUN_ACCESS_KEY_ID 或 ALIYUN_ACCESS_KEY_SECRET 未配置")
	}

	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	u := model.NewArticleRepo(c)
	ctx := svc.NewServiceContext(c, u)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		__.RegisterArticleServiceServer(grpcServer, server.NewArticleServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	serviceGroup.Add(s)

	// Add Kafka consumer
	consumer := mqs.NewArticleConsumer(context.Background(), ctx)
	serviceGroup.Add(kq.MustNewQueue(c.KqConsumerConf, consumer))

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	fmt.Printf("Starting kafka consumer...\n")
	serviceGroup.Start()
}
