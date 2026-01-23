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

type AliGreenConf struct {
	AccessKeyId     string
	AccessKeySecret string
	RegionId        string
	Endpoint        string
}
type Config struct {
	zrpc.RpcServerConf
	Postgres       Postgres
	KqPusherConf   kq.KqConf
	KqConsumerConf kq.KqConf
	AliGreen       AliGreenConf
}
