// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"encoding/json"

	"sea-try-go/service/user/api/internal/svc"
	"sea-try-go/service/user/api/internal/types"
	"sea-try-go/service/user/rpc/pb"

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
	userId := l.ctx.Value("userId").(json.Number)
	Uid, er := userId.Int64()
	if er != nil {
		return nil, er
	}

	rpcReq := &pb.UpdateUserReq{
		Uid:       Uid,
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
		ExtraInfo: req.Extrainfo,
	}

	rpcResp, er := l.svcCtx.UserRpc.UpdateUser(l.ctx, rpcReq)
	if er != nil {
		return nil, er
	}

	return &types.UpdateUserResp{
		User: types.UserInfo{
			Uid:       rpcResp.User.Uid,
			Username:  rpcResp.User.Username,
			Email:     rpcResp.User.Email,
			Extrainfo: rpcResp.User.ExtraInfo,
		},
	}, nil

}
