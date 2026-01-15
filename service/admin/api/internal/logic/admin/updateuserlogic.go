// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"sea-try-go/service/admin/api/internal/svc"
	"sea-try-go/service/admin/api/internal/types"
	"sea-try-go/service/admin/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateuserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateuserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateuserLogic {
	return &UpdateuserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateuserLogic) Updateuser(req *types.UpdateUserReq) (resp *types.UpdateUserResp, err error) {

	rpcReq := &pb.UpdateUserReq{
		Id:        req.Id,
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
		ExtraInfo: req.Extrainfo,
	}

	rpcResp, err := l.svcCtx.AdminRpc.UpdateUser(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.UpdateUserResp{
		Success: true,
		User: types.UserInfo{
			Id:        rpcResp.User.Id,
			Username:  rpcResp.User.Username,
			Email:     rpcResp.User.Email,
			Status:    int64(rpcResp.User.Status),
			Extrainfo: rpcResp.User.ExtraInfo,
		},
	}, nil
}
