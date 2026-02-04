package svc

import (
	"sea-try-go/service/task/rpc/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	pointspb "sea-try-go/service/points/rpc/pb"
)

type ServiceContext struct {
	Config       config.Config
	PointsClient pointspb.PointsServiceClient
	Rdb          *redis.Client
	Gdb          *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.LikeRedis.Addr,
		Password: c.LikeRedis.Pass,
		DB:       c.LikeRedis.DB,
	})

	cli := zrpc.MustNewClient(c.PointsRpc)

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
		Config:       c,
		Rdb:          rdb,
		Gdb:          gdb,
		PointsClient: pointspb.NewPointsServiceClient(cli.Conn()),
	}
}
