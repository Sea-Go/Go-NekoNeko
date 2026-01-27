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

type ResetuserpasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetuserpasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetuserpasswordLogic {
	return &ResetuserpasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetuserpasswordLogic) Resetuserpassword(req *types.ResetUserPasswordReq) (resp *types.ResetUserPasswordResp, err error) {
	rpcReq := &pb.ResetUserPasswordReq{
		Uid: req.Uid,
	}
	rpcResp, err := l.svcCtx.AdminRpc.ResetUserPassword(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return &types.ResetUserPasswordResp{
		Success: rpcResp.Success,
	}, nil
}
