package svc

import (
	"sea-try-go/service/task/rpc/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	Rdb    *redis.Client
	Gdb    *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Pass,
		DB:       c.Redis.DB,
	})

	gdb, err := gorm.Open(postgres.Open(c.Postgres.Dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(c.Postgres.MaxOpenConns)
	sqlDB.SetMaxIdleConns(c.Postgres.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.Postgres.ConnMaxLifetimeMinutes))

	return &ServiceContext{
		Config: c,
		Rdb:    rdb,
		Gdb:    gdb,
	}
}
