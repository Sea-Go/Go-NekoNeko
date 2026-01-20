// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"sea-try-go/service/article/api/internal/config"
	"sea-try-go/service/article/rpc/articleservice"
)

type ServiceContext struct {
	Config     config.Config
	ArticleRpc articleservice.ArticleService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		ArticleRpc: articleservice.NewArticleService(zrpc.MustNewClient(c.ArticleRpcConf)),
	}
}
