package logic

import (
	"context"
	"errors"
	"sea-try-go/service/user/admin/rpc/internal/svc"
	"sea-try-go/service/user/admin/rpc/pb"
	"sea-try-go/service/user/common/cryptx"
	"sea-try-go/service/user/common/errmsg"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pb.LoginReq) (*pb.LoginResp, error) {
	admin, err := l.svcCtx.AdminModel.FindOneAdminByUsername(l.ctx, in.Username)
	if err != nil {
		return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorLoginWrong))
	}
	correct := cryptx.CheckPassword(admin.Password, in.Password)
	if !correct {
		return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorLoginWrong))
	}
	return &pb.LoginResp{
		Uid: admin.Uid,
	}, nil
}
