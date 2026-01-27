// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"errors"
	"sea-try-go/service/user/common/jwt"
	"sea-try-go/service/user/user/api/internal/svc"
	"sea-try-go/service/user/user/api/internal/types"
	"sea-try-go/service/user/user/rpc/pb"
	"time"

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
		Username: req.Username,
		Password: req.Password,
	}

	rpcResp, er := l.svcCtx.UserRpc.Login(l.ctx, rpcReq)

	if er != nil {
		return nil, er
	}

	if rpcResp.Status == 1 {
		return nil, errors.New("用户名或密码错误")
	}

	if rpcResp.Status == 2 {
		return nil, errors.New("用户已被封禁")
	}

	now := time.Now().Unix()
	accessSecret := l.svcCtx.Config.UserAuth.AccessSecret
	accessExpire := l.svcCtx.Config.UserAuth.AccessExpire

	token, e := jwt.GetToken(accessSecret, now, accessExpire, int64(rpcResp.Uid))

	if e != nil {
		return nil, e
	}

	return &types.LoginResp{
		Token: token,
	}, nil
}
