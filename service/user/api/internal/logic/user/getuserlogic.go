// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"encoding/json"
	"errors"

	"sea-try-go/service/user/api/internal/svc"
	"sea-try-go/service/user/api/internal/types"
	"sea-try-go/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetuserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetuserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetuserLogic {
	return &GetuserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetuserLogic) Getuser(req *types.GetUserReq) (resp *types.GetUserResp, err error) {

	userId := l.ctx.Value("userId").(json.Number)
	id, _ := userId.Int64()

	rpcReq := &pb.GetUserReq{
		Id: uint64(id),
	}

	rpcResp, er := l.svcCtx.UserRpc.GetUser(l.ctx, rpcReq)

	if er != nil {
		return nil, er
	}
	if !rpcResp.Found {
		return nil, errors.New("查询个人信息错误")
	}

	return &types.GetUserResp{
		User: types.UserInfo{
			Id:        rpcResp.User.Id,
			Username:  rpcResp.User.Username,
			Email:     rpcResp.User.Email,
			Extrainfo: rpcResp.User.ExtraInfo,
		},
		Found: true,
	}, nil
}
