package logic

import (
	"context"
	"sea-try-go/service/comment/internal/model"
	"time"

	"sea-try-go/service/comment/internal/svc"
	"sea-try-go/service/comment/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentLogic {
	return &GetCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentLogic) GetComment(in *pb.GetCommentReq) (*pb.GetCommentResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn := l.svcCtx.CommentModel
	subject, err := l.svcCtx.CommentCache.GetSubjectWithCache(ctx, in.TargetId, conn)
	if err != nil {
		return nil, err
	}
	var sortType model.ReplySort
	if in.SortType == 1 {
		sortType = model.ReplySortTime
	} else {
		sortType = model.ReplySortTime
	}
	ids, err := l.svcCtx.CommentCache.GetReplyIDsPageCache(ctx, model.GetReplyIDsPageReq{
		TargetType: in.TargetType,
		TargetId:   in.TargetId,
		RootId:     in.RootId,
		Offset:     0,
		Limit:      int(in.Page),
		Sort:       sortType,
		OnlyNormal: false,
	}, conn)
	if err != nil {
		return nil, err
	}
	content, err := l.svcCtx.CommentCache.BatchGetContentCache(ctx, ids, conn)
	if err != nil {
		return nil, err
	}
	return &pb.GetCommentResp{
		Comment: content,
		Subject: subject,
	}, nil
}
