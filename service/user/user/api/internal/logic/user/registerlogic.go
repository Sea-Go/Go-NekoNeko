// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"sea-try-go/service/user/user/api/internal/svc"
	"sea-try-go/service/user/user/api/internal/types"
	"sea-try-go/service/user/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

//register代码存在一定问题:
//用First和Create可能会在高并发情境下产生错误,出现两个人注册了同一个的情况

func (l *RegisterLogic) Register(req *types.CreateUserReq) (resp *types.CreateUserResp, err error) {

	rpcReq := &pb.CreateUserReq{
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
		ExtraInfo: req.Extrainfo,
	}

	rpcResp, e := l.svcCtx.UserRpc.Register(l.ctx, rpcReq)
	//第二个参数必须是指针类型

	if e != nil {
		return nil, e
	}

	return &types.CreateUserResp{
		Uid: rpcResp.Uid,
	}, nil
}
