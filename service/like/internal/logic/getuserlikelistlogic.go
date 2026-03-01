package logic

import (
	"context"

	"sea-try-go/service/like/internal/svc"
	"sea-try-go/service/like/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLikeListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLikeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLikeListLogic {
	return &GetUserLikeListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLikeListLogic) GetUserLikeList(in *pb.GetUserLikeListReq) (*pb.GetUserLikeListResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetUserLikeListResp{}, nil
}
