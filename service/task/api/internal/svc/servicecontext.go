// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"sea-try-go/service/task/api/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
