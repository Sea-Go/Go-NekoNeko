// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"
	"encoding/json"
	"errors"
	"sea-try-go/service/user/admin/api/internal/svc"
	"sea-try-go/service/user/admin/api/internal/types"
	"sea-try-go/service/user/admin/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.UpdateSelfReq) (resp *types.UpdateSelfResp, err error) {
	userId, ok := l.ctx.Value("userId").(json.Number)
	if !ok {
		return nil, errors.New("Token 解析异常")
	}
	uid, err := userId.Int64()
	if err != nil {
		return nil, err
	}
	rpcReq := &pb.UpdateSelfReq{
		Uid:       uid,
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
		ExtraInfo: req.Extrainfo,
	}

	rpcResp, err := l.svcCtx.AdminRpc.UpdateSelf(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.UpdateSelfResp{
		Success: true,
		Admin: types.AdminInfo{
			Uid:       rpcResp.Admin.Uid,
			Username:  rpcResp.Admin.Username,
			Email:     rpcResp.Admin.Email,
			Extrainfo: rpcResp.Admin.ExtraInfo,
		},
	}, nil

}
