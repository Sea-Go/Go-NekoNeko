package svc

import (
	"sea-try-go/service/like/internal/config"
	"sea-try-go/service/like/internal/model"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	BizRedis    *redis.Redis
	DB          *gorm.DB
	KafkaPusher *kq.Pusher
	LikeModel   model.LikeRecordModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	bizRedis := redis.MustNewRedis(c.Storage.Redis.RedisConf)
	dbConn, err := gorm.Open(postgres.Open(c.Storage.Postgres.DataSource), &gorm.Config{})
	if err != nil {
		panic("PostgreSQL连接失败" + err.Error())
	}
	sqlDB, err := dbConn.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
	}
	kafkaPusher := kq.NewPusher(c.KafkaConf.Brokers, c.KafkaConf.Topic)
	return &ServiceContext{
		Config:      c,
		BizRedis:    bizRedis,
		DB:          dbConn,
		KafkaPusher: kafkaPusher,
		LikeModel:   model.NewLikeRecordModel(dbConn),
	}
}
