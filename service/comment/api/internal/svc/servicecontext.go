// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"sea-try-go/service/comment/api/internal/config"
	commentpb "sea-try-go/service/comment/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	CommentCli commentpb.CommentServiceClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := zrpc.MustNewClient(c.CommentRpc)
	return &ServiceContext{
		Config:     c,
		CommentCli: commentpb.NewCommentServiceClient(conn.Conn()),
	}
}
