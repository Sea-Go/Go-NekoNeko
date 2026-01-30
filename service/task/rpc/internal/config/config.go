package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Kafka struct {
		Brokers          []string
		InTopic          string
		OutTopic         string
		GroupKafkaRaw    string
		GroupKafkaFilter string
		GroupGoKa        string
		Offset           string
		Consumers        int
	}
	Redis struct {
		Addr string
		Pass string
		DB   int
	}
	Postgres struct {
		Dsn                    string
		MaxOpenConns           int
		MaxIdleConns           int
		ConnMaxLifetimeMinutes int
	}
}
