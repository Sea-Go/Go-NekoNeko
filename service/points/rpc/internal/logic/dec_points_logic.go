package logic

import (
	"context"

	"sea-try-go/service/points/rpc/internal/svc"
	"sea-try-go/service/points/rpc/pb"

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

func (l *DecPointsLogic) DecPoints(in *__.DecPointsReq) (*__.DecPointsResp, error) {
	// todo: add your logic here and delete this line

	return &__.DecPointsResp{}, nil
}
