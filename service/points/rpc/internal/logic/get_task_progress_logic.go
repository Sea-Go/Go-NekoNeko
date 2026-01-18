package logic

import (
	"context"

	"sea-try-go/service/points/rpc/internal/svc"
	"sea-try-go/service/points/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTaskProgressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTaskProgressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskProgressLogic {
	return &GetTaskProgressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTaskProgressLogic) GetTaskProgress(in *__.GetTaskProgressReq) (*__.GetTaskProgressResp, error) {
	// todo: add your logic here and delete this line

	return &__.GetTaskProgressResp{}, nil
}
