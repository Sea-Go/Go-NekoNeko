package logic

import (
	"context"
	"errors"

	"sea-try-go/service/common/cryptx"
	"sea-try-go/service/common/errmsg"
	"sea-try-go/service/user/rpc/internal/model"
	"sea-try-go/service/user/rpc/internal/svc"
	pb "sea-try-go/service/user/rpc/pb"

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

	user, err := l.svcCtx.UserModel.FindOneByUserName(l.ctx, in.Username)

	if err != nil {
		if err == model.ErrorNotFound {
			return &pb.LoginResp{
				Status: 1,
			}, nil
		}
		return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorDbSelect))
	}

	correct := cryptx.CheckPassword(user.Password, in.Password)
	if !correct {
		return &pb.LoginResp{
			Status: 1,
		}, nil
	}
	if user.Status == 1 {
		return &pb.LoginResp{
			Status: 2,
		}, nil
	}

	return &pb.LoginResp{
		Uid:    user.Uid,
		Status: 0,
	}, nil
}
