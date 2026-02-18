// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"sea-try-go/service/user/admin/api/internal/config"
	"sea-try-go/service/user/admin/rpc/adminservice"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	AdminRpc adminservice.AdminService
}

func NewServiceContext(c config.Config) *ServiceContext {

	return &ServiceContext{
		Config:   c,
		AdminRpc: adminservice.NewAdminService(zrpc.MustNewClient(c.AdminRpc)),
	}
}
