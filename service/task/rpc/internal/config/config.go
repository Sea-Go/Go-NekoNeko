package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	//PointsRpc zrpc.RpcClientConf
	//UserRpc   zrpc.RpcClientConf
	Kafka struct {
		Brokers         []string
		InTopic         string
		RawUserTopic    string
		RawArticleTopic string
		OutUserTopic    string
		OutArticleTopic string

		GroupKafkaRaw    string
		GroupKafkaFilter string

		GroupFilterUser    string
		GroupFilterArticle string

		GroupGoKa        string
		GroupGoKaUser    string
		GroupGoKaArticle string

		Offset    string
		Consumers int
	}
	LikeRedis struct {
		Addr         string
		Pass         string
		DB           int
		ConsumerName string
	}
	Postgres struct {
		Dsn                    string
		MaxOpenConns           int
		MaxIdleConns           int
		ConnMaxLifetimeMinutes int
	}
}
