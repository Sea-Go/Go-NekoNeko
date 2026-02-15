package logic

import (
	"context"

	"sea-try-go/service/follow/rpc/internal/svc"
	"sea-try-go/service/follow/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *FollowLogic) Follow(in *pb.RelationReq) (*pb.BaseResp, error) {
	// 1. 业务特判
	if in.UserId == in.TargetId {
		return &pb.BaseResp{Code: 400, Msg: "Cannot follow yourself"}, nil
	}

	// 2. 呼叫 Model 层干活
	err := l.svcCtx.FollowModel.FollowUser(l.ctx, in.UserId, in.TargetId)
	if err != nil {
		l.Logger.Errorf("FollowUser db err: %v", err)
		return &pb.BaseResp{Code: 500, Msg: "DB Error"}, err
	}

	return &pb.BaseResp{Code: 0, Msg: "success"}, nil
}
