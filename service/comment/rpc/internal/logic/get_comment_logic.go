package logic

import (
	"context"
	"fmt"
	"sea-try-go/service/comment/rpc/internal/model"
	"sea-try-go/service/comment/rpc/internal/svc"
	"sea-try-go/service/comment/rpc/pb"
	"time"

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
	fmt.Println(in.TargetId)
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
	index, err := l.svcCtx.CommentCache.GetCommentIndexCache(ctx, ids, conn)
	if err != nil {
		return nil, err
	}
	content, err := l.svcCtx.CommentCache.BatchGetContentCache(ctx, ids, conn)
	if err != nil {
		return nil, err
	}
	comment := make([]*pb.CommentItem, 0, len(content))
	for i, _ := range index {
		u := index[i]
		v := content[i]
		comment = append(comment, &pb.CommentItem{
			Id:           u.Id,
			UserId:       u.UserId,
			Content:      v.Content,
			RootId:       u.RootId,
			ParentId:     u.ParentId,
			LikeCount:    u.LikeCount,
			DislikeCount: u.DislikeCount,
			ReplyCount:   u.ReplyCount,
			Attribute:    u.Attribute,
			State:        pb.CommentState(u.State),
			CreatedAt:    u.CreatedAt.Format("2006-01-02 15:04:05"),
			Meta:         v.Meta,
			Children:     nil, //日后再说
		})
	}
	return &pb.GetCommentResp{
		Comment: comment,
		Subject: &pb.SubjectInfo{
			TargetType: subject.TargetType,
			TargetId:   subject.TargetId,
			TotalCount: subject.TotalCount,
			RootCount:  subject.RootCount,
			State:      pb.SubjectState(subject.State),
			Attribute:  subject.Attribute,
		},
	}, nil
}
