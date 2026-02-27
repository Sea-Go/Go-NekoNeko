package svc

import (
	"sea-try-go/service/comment/internal/cache"
	"sea-try-go/service/comment/internal/config"
	"sea-try-go/service/comment/internal/model"
	"sea-try-go/service/comment/internal/utils"

	"github.com/zeromicro/go-queue/kq"
)

type ServiceContext struct {
	Config          config.Config
	CommentModel    *model.CommentModel
	CommentCache    *cache.CommentCache
	KqPusherClient  *kq.Pusher
	SensitiveFilter *utils.SensitiveFilter
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := model.InitDB(c.DB.DataSource)
	rdb := cache.InitRedis(c.BizRedis.Host)
	pusher := kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic)
	blackwords := []string{"傻逼", "我草"}
	return &ServiceContext{
		Config:          c,
		CommentModel:    model.NewCommentModel(db),
		CommentCache:    cache.NewCommentCache(rdb),
		KqPusherClient:  pusher,
		SensitiveFilter: utils.NewSensitiveFilter(blackwords),
	}
}
