package logic

import (
	"context"

	"sea-try-go/service/follow/rpc/internal/svc"
	"sea-try-go/service/follow/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnblockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnblockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnblockLogic {
	return &UnblockLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UnblockLogic) Unblock(in *pb.RelationReq) (*pb.BaseResp, error) {
	err := l.svcCtx.FollowModel.UnblockUser(l.ctx, in.UserId, in.TargetId)
	if err != nil {
		l.Logger.Errorf("UnblockUser db err: %v", err)
		return &pb.BaseResp{Code: 500, Msg: "DB Error"}, err
	}
	return &pb.BaseResp{Code: 0, Msg: "success"}, nil
}
