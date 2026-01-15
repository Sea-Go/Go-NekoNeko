package logic

import (
	"context"

	"sea-try-go/service/user/rpc/internal/model"
	"sea-try-go/service/user/rpc/internal/svc"
	pb "sea-try-go/service/user/rpc/pb"

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

	err := l.svcCtx.DB.Where("id = ?", in.Id).Delete(&model.User{}).Error
	if err != nil {
		return &pb.DeleteUserResp{
			Success: false,
		}, err
	}

	return &pb.DeleteUserResp{
		Success: true,
	}, nil
}
