package logic

import (
	"context"

	"sea-try-go/service/like/internal/svc"
	"sea-try-go/service/like/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLikeStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLikeStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLikeStateLogic {
	return &GetLikeStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLikeStateLogic) GetLikeState(in *pb.GetLikeStateReq) (*pb.GetLikeStateResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetLikeStateResp{}, nil
}
