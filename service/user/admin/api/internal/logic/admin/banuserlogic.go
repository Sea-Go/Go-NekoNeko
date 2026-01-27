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

type BanuserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBanuserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BanuserLogic {
	return &BanuserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BanuserLogic) Banuser(req *types.BanUserReq) (resp *types.BanUserResp, err error) {
	uid := req.Uid
	rpcReq := &pb.BanUserReq{
		Uid: uid,
	}
	rpcResp, err := l.svcCtx.AdminRpc.BanUser(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.BanUserResp{
		Success: rpcResp.Success,
	}, nil
}
