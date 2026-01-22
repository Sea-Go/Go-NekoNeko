package logic

import (
	"context"
	"errors"

	"sea-try-go/service/admin/rpc/internal/model"
	"sea-try-go/service/admin/rpc/internal/svc"
	"sea-try-go/service/admin/rpc/pb"
	"sea-try-go/service/common/errmsg"

	"github.com/zeromicro/go-zero/core/logx"
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
			return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorUserNotExist))
		}
		return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorServerCommon))
	}
	return &pb.UnBanUserResp{
		Success: true,
	}, nil
}
