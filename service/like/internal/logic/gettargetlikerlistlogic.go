package logic

import (
	"context"

	"sea-try-go/service/like/internal/svc"
	"sea-try-go/service/like/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTargetLikerListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTargetLikerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTargetLikerListLogic {
	return &GetTargetLikerListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTargetLikerListLogic) GetTargetLikerList(in *pb.GetTargetLikerListReq) (*pb.GetTargetLikerListResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetTargetLikerListResp{}, nil
}
