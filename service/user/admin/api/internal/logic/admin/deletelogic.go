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

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.DeleteUserReq) (resp *types.DeleteUserResp, err error) {
	rpcReq := &pb.DeleteUserReq{
		Uid: req.Uid,
	}
	rpcResp, err := l.svcCtx.AdminRpc.DeleteUser(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return &types.DeleteUserResp{
		Success: rpcResp.Success,
	}, nil
}
