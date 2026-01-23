package logic

import (
	"context"

	"sea-try-go/service/hotspot/rpc/internal/svc"
	"sea-try-go/service/hotspot/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCommentsLogic {
	return &ListCommentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取评论列表
func (l *ListCommentsLogic) ListComments(in *__.ListCommentReq) (*__.ListCommentResp, error) {
	// todo: add your logic here and delete this line

	return &__.ListCommentResp{}, nil
}
