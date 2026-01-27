package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Kafka struct {
		Brokers   []string
		Topic     string
		Group     string
		Offset    string
		Consumers int
	}
}
