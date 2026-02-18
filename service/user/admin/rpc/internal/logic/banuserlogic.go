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

type BanUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBanUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BanUserLogic {
	return &BanUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BanUserLogic) BanUser(in *pb.BanUserReq) (*pb.BanUserResp, error) {
	err := l.svcCtx.AdminModel.UpdateUserStatusByUid(l.ctx, in.Uid, 1)
	if err != nil {
		if err == model.ErrorNotFound {
			logger.LogBusinessErr(l.ctx, errmsg.ErrorUserNotExist, err)
			return nil, status.Error(codes.NotFound, "用户不存在")
		}
		logger.LogBusinessErr(l.ctx, errmsg.ErrorDbUpdate, err)
		return nil, status.Error(codes.Internal, "DB更新失败")
	}
	logger.LogInfo(l.ctx, fmt.Sprintf("ban user success,uid : %d", in.Uid))
	return &pb.BanUserResp{
		Success: true,
	}, nil
}
