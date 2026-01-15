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

type DeleteuserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteuserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteuserLogic {
	return &DeleteuserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteuserLogic) Deleteuser(req *types.DeleteUserReq) (resp *types.DeleteUserResp, err error) {

	userId := l.ctx.Value("userId").(json.Number)
	id, _ := userId.Int64()

	rpcReq := &pb.DeleteUserReq{
		Id: uint64(id),
	}

	rpcResp, er := l.svcCtx.UserRpc.DeleteUser(l.ctx, rpcReq)
	if er != nil {
		return nil, er
	}
	if !rpcResp.Success {
		return &types.DeleteUserResp{
			Success: false,
		}, nil
	}

	return &types.DeleteUserResp{
		Success: true,
	}, nil
}
