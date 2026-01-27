// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"
	"sea-try-go/service/user/admin/api/internal/svc"
	"sea-try-go/service/user/admin/api/internal/types"
	"sea-try-go/service/user/admin/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnbanuserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnbanuserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbanuserLogic {
	return &UnbanuserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnbanuserLogic) Unbanuser(req *types.UnBanUserReq) (resp *types.UnBanUserResp, err error) {
	rpcReq := &pb.UnBanUserReq{
		Uid: req.Uid,
	}
	rpcResp, er := l.svcCtx.AdminRpc.UnBanUser(l.ctx, rpcReq)
	if er != nil {
		return nil, er
	}
	return &types.UnBanUserResp{
		Success: rpcResp.Success,
	}, nil
}
