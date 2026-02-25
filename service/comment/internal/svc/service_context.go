package svc

import (
	"sea-try-go/service/comment/internal/config"
	"sea-try-go/service/comment/internal/model"
	"sea-try-go/service/comment/internal/utils"

	"github.com/zeromicro/go-queue/kq"
)

type ServiceContext struct {
	Config          config.Config
	CommentModel    *model.CommentModel
	KqPusherClient  *kq.Pusher
	SensitiveFilter *utils.SensitiveFilter
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := model.InitDB(c.DB.DataSource)
	pusher := kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic)
	blackwords := []string{"傻逼", "我草"}
	return &ServiceContext{
		Config:          c,
		CommentModel:    model.NewCommentModel(db),
		KqPusherClient:  pusher,
		SensitiveFilter: utils.NewSensitiveFilter(blackwords),
	}
}
