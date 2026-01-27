// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"
	"time"

	"sea-try-go/service/user/admin/api/internal/svc"
	"sea-try-go/service/user/admin/api/internal/types"
	"sea-try-go/service/user/admin/rpc/pb"
	"sea-try-go/service/user/common/jwt"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	rpcReq := &pb.LoginReq{
		Password: req.Password,
		Username: req.Username,
	}

	rpcResp, err := l.svcCtx.AdminRpc.Login(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	accessSecret := l.svcCtx.Config.AdminAuth.AccessSecret
	accessExpire := l.svcCtx.Config.AdminAuth.AccessExpire
	token, err := jwt.GetToken(accessSecret, now, accessExpire, int64(rpcResp.Uid))
	if err != nil {
		return nil, err
	}
	return &types.LoginResp{
		Token: token,
	}, nil
}
