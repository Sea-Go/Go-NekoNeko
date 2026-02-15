package logic

import (
	"context"

	"sea-try-go/service/follow/rpc/internal/svc"
	"sea-try-go/service/follow/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowerListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowerListLogic {
	return &GetFollowerListLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetFollowerListLogic) GetFollowerList(in *pb.ListReq) (*pb.UserListResp, error) {
	ids, err := l.svcCtx.FollowModel.GetFollowerList(l.ctx, in.UserId, in.Offset, in.Limit)
	if err != nil {
		l.Logger.Errorf("GetFollowerList db err: %v", err)
		return &pb.UserListResp{Code: 500, Msg: "DB Error"}, err
	}
	return &pb.UserListResp{Code: 0, Msg: "success", UserIds: ids}, nil
}
