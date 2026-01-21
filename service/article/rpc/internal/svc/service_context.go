package svc

import (
	"github.com/zeromicro/go-queue/kq"
	"sea-try-go/common/utils"
	"sea-try-go/service/article/rpc/internal/config"
	"sea-try-go/service/article/rpc/internal/model"
)

type ServiceContext struct {
	Config      config.Config
	ArticleRepo *model.ArticleRepo
	KqPusher    *kq.Pusher
}

func NewServiceContext(c config.Config, articleRepo *model.ArticleRepo) *ServiceContext {
	snowflake.Init()
	return &ServiceContext{
		Config:      c,
		ArticleRepo: articleRepo,
		KqPusher:    kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
	}
}
