// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"sea-try-go/service/user/user/api/internal/config"
	"sea-try-go/service/user/user/rpc/userservice"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc userservice.UserService
	Pg      *pgxpool.Pool
}

func NewServiceContext(c config.Config, pgPool *pgxpool.Pool) *ServiceContext {

	return &ServiceContext{
		Config:  c,
		UserRpc: userservice.NewUserService(zrpc.MustNewClient(c.UserRpc)),
		Pg:      pgPool,
	}
}
