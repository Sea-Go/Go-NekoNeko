package logic

import (
	"context"

	"sea-try-go/service/user/rpc/internal/svc"
	"sea-try-go/service/user/rpc/pb"

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
	err := l.svcCtx.UserModel.Delete(in.Uid)
	if err != nil {
		l.Errorf("删除用户失败: %v", err)
		return &pb.DeleteUserResp{Success: false}, err
	}

	return &pb.DeleteUserResp{Success: true}, nil
}
