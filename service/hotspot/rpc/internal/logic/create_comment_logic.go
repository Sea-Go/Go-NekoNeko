package logic

import (
	"context"

	"sea-try-go/service/hotspot/rpc/internal/svc"
	"sea-try-go/service/hotspot/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 发表评论
func (l *CreateCommentLogic) CreateComment(in *__.CommentReq) (*__.CommentResp, error) {
	// todo: add your logic here and delete this line

	return &__.CommentResp{}, nil
}
