package svc

import (
	cache2 "sea-try-go/service/comment/rpc/internal/cache"
	"sea-try-go/service/comment/rpc/internal/config"
	model2 "sea-try-go/service/comment/rpc/internal/model"
	"sea-try-go/service/comment/rpc/internal/utils"

	"github.com/zeromicro/go-queue/kq"
)

type ServiceContext struct {
	Config          config.Config
	CommentModel    *model2.CommentModel
	CommentCache    *cache2.CommentCache
	KqPusherClient  *kq.Pusher
	SensitiveFilter *utils.SensitiveFilter
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := model2.InitDB(c.DB.DataSource)
	rdb := cache2.InitRedis(c.BizRedis.Host)
	pusher := kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic)
	blackwords := []string{"傻逼", "我草"}
	return &ServiceContext{
		Config:          c,
		CommentModel:    model2.NewCommentModel(db),
		CommentCache:    cache2.NewCommentCache(rdb),
		KqPusherClient:  pusher,
		SensitiveFilter: utils.NewSensitiveFilter(blackwords),
	}
}
