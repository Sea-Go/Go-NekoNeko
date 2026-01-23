package logic

import (
	"context"

	"sea-try-go/service/hotspot/rpc/internal/svc"
	"sea-try-go/service/hotspot/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsLikedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsLikedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsLikedLogic {
	return &IsLikedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询是否点赞过
func (l *IsLikedLogic) IsLiked(in *__.IsLikedReq) (*__.IsLikedResp, error) {
	// todo: add your logic here and delete this line

	return &__.IsLikedResp{}, nil
}
