package svc

import (
	"sea-try-go/service/points/rpc/internal/config"
	"sea-try-go/service/points/rpc/internal/model"
	"time"

	"github.com/IBM/sarama"
	"github.com/RussellLuo/timingwheel"
	"github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	Config        config.Config
	PointsModel   *model.PointsModel
	RDB           *redis.Client
	KafKa         sarama.ConsumerGroup
	KafKaProducer sarama.SyncProducer
	TimingWheel   *timingwheel.TimingWheel
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
		Addr:     c.RedisConf.Addr,
		Password: c.RedisConf.Password,
		DB:       c.RedisConf.DB,
		PoolSize: c.RedisConf.PoolSize,
	})
	kafka := sarama.NewConfig()
	kafka.Producer.RequiredAcks = sarama.WaitForAll
	kafka.Producer.Retry.Max = 5
	kafka.Producer.Return.Successes = true
	kafka.Producer.Idempotent = true
	kafka.Net.MaxOpenRequests = 1
	proConfig := sarama.NewConfig()
	proConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	proConfig.Consumer.Return.Errors = true
	proConfig.Producer.Return.Successes = true
	group, err := sarama.NewConsumerGroup(c.KafkaConf.Brokers, c.KafkaConf.Group, proConfig)
	if err != nil {
		panic("kafka1初始化失败" + err.Error())
	}
	producer, err := sarama.NewSyncProducer(c.KafkaConf.Brokers, kafka)
	if err != nil {
		panic("kafka初始化失败" + err.Error())
	}
	if err != nil {
		panic("kafka连接失败:" + err.Error())
	}
	tw := timingwheel.NewTimingWheel(time.Second, 60)
	tw.Start()
	return &ServiceContext{
		Config:        c,
		PointsModel:   model.NewPointsModel(db),
		RDB:           rdb,
		KafKa:         group,
		KafKaProducer: producer,
	}
}
