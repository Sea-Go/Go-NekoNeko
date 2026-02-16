package logic

import (
	"context"
	"fmt"
	"sea-try-go/service/common/logger"
	"sea-try-go/service/user/common/cryptx"
	"sea-try-go/service/user/common/errmsg"
	"sea-try-go/service/user/user/rpc/internal/model"
	"sea-try-go/service/user/user/rpc/internal/svc"
	"sea-try-go/service/user/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
			logger.LogInfo(l.ctx, fmt.Sprintf("login failed : user not found,username : %s", in.Username))
			return &pb.LoginResp{
				Status: 1,
			}, nil
		}
		logger.LogBusinessErr(l.ctx, errmsg.ErrorDbSelect, err)
		return nil, status.Error(codes.Internal, "DB查询失败")
	}

	correct := cryptx.CheckPassword(user.Password, in.Password)
	if !correct {
		logger.LogInfo(l.ctx, fmt.Sprintf("login failed: username or password incorrect,username : %s", in.Username))
		return &pb.LoginResp{
			Status: 1,
		}, nil
	}
	if user.Status == 1 {
		logger.LogInfo(l.ctx, fmt.Sprintf("login failed: user banned, username :  %s", in.Username))
		return &pb.LoginResp{
			Status: 2,
		}, nil
	}
	logger.LogInfo(l.ctx, fmt.Sprintf("login success,username : %s", in.Username))
	return &pb.LoginResp{
		Uid:    user.Uid,
		Status: 0,
	}, nil
}
