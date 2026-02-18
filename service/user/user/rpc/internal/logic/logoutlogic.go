package logic

import (
	"context"
	"sea-try-go/service/user/user/rpc/internal/svc"
	"sea-try-go/service/user/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LogoutLogic) Logout(in *pb.LogoutReq) (*pb.LogoutResp, error) {
	return &pb.LogoutResp{}, nil
}
