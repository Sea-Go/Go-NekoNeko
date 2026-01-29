package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Kafka struct {
		Brokers    []string
		InTopic    string
		OutTopic   string
		GroupKafka string
		GroupGoKa  string
		Offset     string
		Consumers  int
	}
}
