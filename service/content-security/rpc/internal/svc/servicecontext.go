package svc

import (
	"sea-try-go/service/common/logger"
	"sea-try-go/service/content-security/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	logx.MustSetup(c.LogConf)
	logger.Init(c.Name)
	return &ServiceContext{
		Config: c,
	}
}
