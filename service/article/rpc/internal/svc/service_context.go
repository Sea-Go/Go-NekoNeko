package svc

import (
	"log"
	"sea-try-go/service/article/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config   config.Config
	RedisCli *redis.Redis
	DB       *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(postgres.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		log.Fatalln("数据库连接失败")
	}
	return &ServiceContext{
		Config:   c,
		DB:       db,
		RedisCli: redis.MustNewRedis(c.RedisCli), // MustNewRedis 连接失败会直接 panic
	}
}
