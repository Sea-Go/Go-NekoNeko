package logic

import (
	"context"
	"fmt"
	"sea-try-go/service/common/logger"
	"sea-try-go/service/user/admin/rpc/internal/model"
	"sea-try-go/service/user/admin/rpc/internal/svc"
	"sea-try-go/service/user/admin/rpc/pb"
	"sea-try-go/service/user/common/errmsg"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UnBanUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnBanUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnBanUserLogic {
	return &UnBanUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnBanUserLogic) UnBanUser(in *pb.UnBanUserReq) (*pb.UnBanUserResp, error) {
	err := l.svcCtx.AdminModel.UpdateUserStatusByUid(l.ctx, in.Uid, 0)
	if err != nil {
		if err == model.ErrorNotFound {
			logger.LogBusinessErr(l.ctx, errmsg.ErrorUserNotExist, err)
			return nil, status.Error(codes.NotFound, "用户不存在")
		}
		logger.LogBusinessErr(l.ctx, errmsg.ErrorDbUpdate, err)
		return nil, status.Error(codes.Internal, "DB更新失败")
	}
	logger.LogInfo(l.ctx, fmt.Sprintf("unban user success,uid : %d", in.Uid))
	return &pb.UnBanUserResp{
		Success: true,
	}, nil
}
