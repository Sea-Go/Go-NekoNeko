package logic

import (
	"context"
	"fmt"

	"sea-try-go/service/comment/internal/svc"
	"sea-try-go/service/comment/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeCommentLogic {
	return &LikeCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LikeCommentLogic) LikeComment(in *pb.LikeCommentReq) (*pb.LikeCommentResp, error) {
	if in.CommentId == 0 {
		return nil, fmt.Errorf("参数错误:评论ID不能为空")
	}
	ownerId, err := l.svcCtx.CommentModel.GetOwnerId(l.ctx, in.TargetType, in.TargetId)
	if err != nil {
		return nil, fmt.Errorf("查询Owner失败")
	}
	err = l.svcCtx.CommentModel.LikeCommentTx(l.ctx, in.UserId, in.CommentId, in.TargetType, in.TargetId, in.ActionType, int64(ownerId))
	if err != nil {
		return nil, err
	}

	return &pb.LikeCommentResp{
		Success: true,
	}, nil
}
