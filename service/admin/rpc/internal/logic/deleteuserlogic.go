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
	err := l.svcCtx.AdminModel.DeleteOneUserByUid(l.ctx, in.Uid)
	if err != nil {
		if err == model.ErrorNotFound {
			return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorUserNotExist))
		}
	}
	return &pb.DeleteUserResp{
		Success: true,
	}, nil
}
