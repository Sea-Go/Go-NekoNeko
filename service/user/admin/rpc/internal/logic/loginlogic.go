package logic

import (
	"context"
	"fmt"
	"sea-try-go/service/common/logger"
	"sea-try-go/service/user/admin/rpc/internal/model"
	"sea-try-go/service/user/admin/rpc/internal/svc"
	"sea-try-go/service/user/admin/rpc/pb"
	"sea-try-go/service/user/common/cryptx"
	"sea-try-go/service/user/common/errmsg"

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
	admin, err := l.svcCtx.AdminModel.FindOneAdminByUsername(l.ctx, in.Username)
	if err != nil {
		if err == model.ErrorNotFound {
			logger.LogBusinessErr(l.ctx, errmsg.ErrorUserNotExist, err)
			return nil, status.Error(codes.NotFound, "用户不存在")
		}
		logger.LogBusinessErr(l.ctx, errmsg.ErrorDbSelect, err)
		return nil, status.Error(codes.Internal, "DB查询失败")
	}
	correct := cryptx.CheckPassword(admin.Password, in.Password)
	if !correct {
		logger.LogBusinessErr(l.ctx, errmsg.ErrorLoginWrong, fmt.Errorf("password mismatched"))
		return nil, status.Error(codes.Unauthenticated, "密码错误")
	}
	logger.LogInfo(l.ctx, fmt.Sprintf("login success,username:%s", in.Username))
	return &pb.LoginResp{
		Uid: admin.Uid,
	}, nil
}
