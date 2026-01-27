// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"encoding/json"
	"sea-try-go/service/user/user/api/internal/svc"
	"sea-try-go/service/user/user/api/internal/types"
	"sea-try-go/service/user/user/rpc/pb"

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
	uid, _ := userId.Int64()

	rpcReq := &pb.DeleteUserReq{
		Uid: int64(uid),
	}

	rpcResp, err := l.svcCtx.UserRpc.DeleteUser(l.ctx, rpcReq)
	if err != nil {
		return nil, err
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
