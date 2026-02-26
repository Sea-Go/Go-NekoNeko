package logic

import (
	"context"
	"fmt"

	"sea-try-go/service/comment/internal/svc"
	"sea-try-go/service/comment/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCommentLogic) DeleteComment(in *pb.DeleteCommentReq) (*pb.DeleteCommentResp, error) {
	if in.CommentId == 0 || in.TargetId == "" {
		return nil, fmt.Errorf("参数错误:评论ID和目标ID不能为空")
	}
	remainCount, err := l.svcCtx.CommentModel.DeleteCommentTx(l.ctx, in.CommentId, in.UserId, in.TargetType, in.TargetId)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteCommentResp{
		Success:           true,
		SubjectTotalCount: remainCount,
	}, nil
}
