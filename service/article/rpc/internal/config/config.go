package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/zrpc"
)

type Postgres struct {
	Host     string
	Dbname   string
	Password string
	Port     string
	User     string
}
type Config struct {
	zrpc.RpcServerConf
	Postgres     Postgres
	KqPusherConf kq.KqConf
}
