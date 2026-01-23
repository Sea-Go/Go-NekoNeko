package svc

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	green "github.com/alibabacloud-go/green-20220302/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"sea-try-go/common/utils"
	"sea-try-go/service/article/rpc/internal/config"
	"sea-try-go/service/article/rpc/internal/model"
)

type ServiceContext struct {
	Config      config.Config
	ArticleRepo *model.ArticleRepo
	KqPusher    *kq.Pusher
	GreenClient *green.Client
}

func NewServiceContext(c config.Config, articleRepo *model.ArticleRepo) *ServiceContext {
	snowflake.Init()
	config := &openapi.Config{
		AccessKeyId:     &c.AliGreen.AccessKeyId,
		AccessKeySecret: &c.AliGreen.AccessKeySecret,
		Endpoint:        tea.String(c.AliGreen.Endpoint),
		ConnectTimeout:  tea.Int(3000),
		ReadTimeout:     tea.Int(6000),
	}
	client, err := green.NewClient(config)
	if err != nil {
		logx.Errorf("Failed to init AliGreen client: %v", err)
	}
	return &ServiceContext{
		Config:      c,
		ArticleRepo: articleRepo,
		KqPusher:    kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
		GreenClient: client,
	}
}
