package svc

import (
	"sea-try-go/service/article/rpc/internal/config"
	model "sea-try-go/service/article/rpc/internal/model/postgres"
	snowflake "sea-try-go/service/common/utils"
)

type ServiceContext struct {
	Config      config.Config
	ArticleRepo *model.ArticleRepo
}

func NewServiceContext(c config.Config, articleRepo *model.ArticleRepo) *ServiceContext {
	snowflake.Init()
	return &ServiceContext{
		Config:      c,
		ArticleRepo: articleRepo,
	}
}
