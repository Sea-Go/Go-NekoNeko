package logic

import (
	"context"

	"sea-try-go/service/follow/rpc/internal/svc"
	"sea-try-go/service/follow/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnfollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnfollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnfollowLogic {
	return &UnfollowLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UnfollowLogic) Unfollow(in *pb.RelationReq) (*pb.BaseResp, error) {
	err := l.svcCtx.FollowModel.UnfollowUser(l.ctx, in.UserId, in.TargetId)
	if err != nil {
		l.Logger.Errorf("UnfollowUser db err: %v", err)
		return &pb.BaseResp{Code: 500, Msg: "DB Error"}, err
	}
	return &pb.BaseResp{Code: 0, Msg: "success"}, nil
}
