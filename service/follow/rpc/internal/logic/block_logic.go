package logic

import (
	"context"

	"sea-try-go/service/follow/rpc/internal/svc"
	"sea-try-go/service/follow/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlockLogic {
	return &BlockLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *BlockLogic) Block(in *pb.RelationReq) (*pb.BaseResp, error) {
	if in.UserId == in.TargetId {
		return &pb.BaseResp{Code: 400, Msg: "Cannot block yourself"}, nil
	}

	err := l.svcCtx.FollowModel.BlockUser(l.ctx, in.UserId, in.TargetId)
	if err != nil {
		l.Logger.Errorf("BlockUser db err: %v", err)
		return &pb.BaseResp{Code: 500, Msg: "DB Error"}, err
	}
	return &pb.BaseResp{Code: 0, Msg: "success"}, nil
}
