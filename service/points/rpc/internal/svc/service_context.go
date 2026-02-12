package svc

import (
	"sea-try-go/service/points/rpc/internal/config"
	"sea-try-go/service/points/rpc/internal/model"
	"sea-try-go/service/user/rpc/userservice"

	"github.com/zeromicro/go-queue/dq"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config              config.Config
	PointsModel         *model.PointsModel
	UserRpc             userservice.UserService
	RetryDqPusherClient dq.Producer
	RetryDqConsumer     dq.Consumer
	DqPusherClient      dq.Producer
	DqConsumer          dq.Consumer
	KqPusherClient      *kq.Pusher
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

	return &ServiceContext{
		Config:              c,
		PointsModel:         model.NewPointsModel(db),
		RetryDqPusherClient: dq.NewProducer(c.RetryDqConf.Beanstalks),
		RetryDqConsumer:     dq.NewConsumer(c.RetryDqConf),
		DqPusherClient:      dq.NewProducer(c.DqConf.Beanstalks),
		DqConsumer:          dq.NewConsumer(c.DqConf),
		UserRpc:             userservice.NewUserService(zrpc.MustNewClient(c.UserRpcConf)),
		KqPusherClient:      kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
	}
}
