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

type ResetUserPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetUserPasswordLogic {
	return &ResetUserPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResetUserPasswordLogic) ResetUserPassword(in *pb.ResetUserPasswordReq) (*pb.ResetUserPasswordResp, error) {

	_, err := l.svcCtx.AdminModel.FindOneUserByUid(l.ctx, in.Uid)
	if err != nil {
		if err == model.ErrorNotFound {
			logger.LogBusinessErr(l.ctx, errmsg.ErrorUserNotExist, err)
			return nil, status.Error(codes.NotFound, "用户不存在")
		}
		logger.LogBusinessErr(l.ctx, errmsg.ErrorDbSelect, err)
		return nil, status.Error(codes.Internal, "DB查询失败")
	}
	var password string
	password, err = cryptx.PasswordEncrypt(l.svcCtx.Config.System.DefaultPassword)
	if err != nil {
		logger.LogBusinessErr(l.ctx, errmsg.ErrorServerCommon, err)
		return nil, status.Error(codes.Internal, "密码加密失败")
	}
	err = l.svcCtx.AdminModel.UpdateUserPasswordByUid(l.ctx, in.Uid, password)
	if err != nil {
		logger.LogBusinessErr(l.ctx, errmsg.ErrorDbUpdate, err)
		return nil, status.Error(codes.Internal, "DB更新失败")
	}
	logger.LogInfo(l.ctx, fmt.Sprintf("reset password success,uid : %d", in.Uid))
	return &pb.ResetUserPasswordResp{
		Success: true,
	}, nil
}
