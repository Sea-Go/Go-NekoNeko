package logic

import (
	"context"
	"sea-try-go/service/common/logger"
	"sea-try-go/service/user/common/errmsg"
	"sea-try-go/service/user/user/rpc/internal/model"
	"sea-try-go/service/user/user/rpc/internal/svc"
	"sea-try-go/service/user/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeleteUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteUserLogic) DeleteUser(in *pb.DeleteUserReq) (*pb.DeleteUserResp, error) {

	_, err := l.svcCtx.UserModel.FindOneByUid(l.ctx, in.Uid)
	if err != nil {
		if err == model.ErrorNotFound {
			logger.LogBusinessErr(l.ctx, errmsg.ErrorUserNotExist, err)
			return &pb.DeleteUserResp{
				Success: false,
			}, status.Error(codes.NotFound, "用户不存在")
		}
		logger.LogBusinessErr(l.ctx, errmsg.ErrorDbSelect, err)
		return &pb.DeleteUserResp{
			Success: false,
		}, status.Error(codes.Internal, "DB查询失败")
	}
	err = l.svcCtx.UserModel.DeleteUserByUid(l.ctx, in.Uid)
	if err != nil {
		logger.LogBusinessErr(l.ctx, errmsg.ErrorDbUpdate, err)
		return &pb.DeleteUserResp{
			Success: false,
		}, status.Error(codes.Internal, "DB更新失败")
	}
	logger.LogInfo(l.ctx, "delete success")
	return &pb.DeleteUserResp{
		Success: true,
	}, nil
}
