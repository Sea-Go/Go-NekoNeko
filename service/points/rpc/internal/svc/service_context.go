package svc

import (
	"sea-try-go/service/points/rpc/internal/config"
	"sea-try-go/service/points/rpc/internal/model"

	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	Config      config.Config
	PointsModel *model.PointsModel
	RDB         *redis.Client
	KafKa       sarama.SyncProducer
}

func NewServiceContext(c config.Config) *ServiceContext {
	dbConfig := model.DBConf{
		Host:     c.Postgres.Host,
		Port:     c.Postgres.Port,
		User:     c.Postgres.User,
		Password: c.Postgres.Password,
		DBName:   c.Postgres.DBName,
		Mode:     c.Postgres.Mode,
	}
	db := model.InitDB(dbConfig)
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
		PoolSize: c.Redis.PoolSize,
	})
	kafka := sarama.NewConfig()
	kafka.Producer.RequiredAcks = sarama.WaitForAll
	kafka.Producer.Retry.Max = 5
	kafka.Producer.Return.Successes = true
	kafka.Producer.Idempotent = true
	kfk, err := sarama.NewSyncProducer(c.KafkaConf.Brokers, kafka)
	if err != nil {
		panic("kafka连接失败:" + err.Error())
	}
	return &ServiceContext{
		Config:      c,
		PointsModel: model.NewPointsModel(db),
		RDB:         rdb,
		KafKa:       kfk,
	}
}
