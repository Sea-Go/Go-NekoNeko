package logic

import (
	"context"

	"sea-try-go/service/points/rpc/internal/svc"
	pb "sea-try-go/service/points/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddPointsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddPointsLogic {
	return &AddPointsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddPointsLogic) AddPoints(in *pb.AddPointsReq) (*pb.AddPointsResp, error) {
	result, err := ProcessPoints(l.ctx, l.svcCtx, in.UserId, in.RequestId, in.AddPoints)
	if err != nil {
		return nil, err
	}
	return &pb.AddPointsResp{Success: result.Success, Message: result.Message}, nil
}
