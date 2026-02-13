package logic

import (
	"context"

	"sea-try-go/service/points/rpc/internal/svc"
	pb "sea-try-go/service/points/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DecPointsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDecPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DecPointsLogic {
	return &DecPointsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DecPointsLogic) DecPoints(in *pb.DecPointsReq) (*pb.DecPointsResp, error) {
	// 扣积分使用负数
	result, err := ProcessPoints(l.ctx, l.svcCtx, in.UserId, in.RequestId, -in.DecPoints)
	if err != nil {
		return nil, err
	}
	return &pb.DecPointsResp{Success: result.Success, Message: result.Message}, nil
}
