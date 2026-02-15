package logic

import (
	"context"

	"sea-try-go/service/follow/rpc/internal/svc"
	"sea-try-go/service/follow/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowListLogic {
	return &GetFollowListLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetFollowListLogic) GetFollowList(in *pb.ListReq) (*pb.UserListResp, error) {
	ids, err := l.svcCtx.FollowModel.GetFollowList(l.ctx, in.UserId, in.Offset, in.Limit)
	if err != nil {
		l.Logger.Errorf("GetFollowList db err: %v", err)
		return &pb.UserListResp{Code: 500, Msg: "DB Error"}, err
	}
	return &pb.UserListResp{Code: 0, Msg: "success", UserIds: ids}, nil
}
