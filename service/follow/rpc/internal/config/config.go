package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Neo4j struct {
		Uri      string
		Username string
		Password string
	}
}
