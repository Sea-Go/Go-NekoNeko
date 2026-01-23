package config

import "github.com/zeromicro/go-zero/zrpc"

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
	Redis struct {
		Addr     string
		Password string
		DB       int
		PoolSize int
	}
	KafkaConf struct {
		Brokers []string
		Topic   string
		Group   string
	}
}
