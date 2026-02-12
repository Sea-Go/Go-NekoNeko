package config

import (
	"github.com/zeromicro/go-queue/dq"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Postgres struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
		Mode     string
	}
	UserRpcConf  zrpc.RpcClientConf
	RetryDqConf  dq.DqConf
	DqConf       dq.DqConf
	KqPusherConf struct {
		Brokers []string
		Topic   string
	}
	KqConsumerConf kq.KqConf
}
