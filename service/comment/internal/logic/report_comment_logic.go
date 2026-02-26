package logic

import (
	"context"
	"fmt"

	"sea-try-go/service/comment/internal/model"
	"sea-try-go/service/comment/internal/svc"
	"sea-try-go/service/comment/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReportCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReportCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportCommentLogic {
	return &ReportCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReportCommentLogic) ReportComment(in *pb.ReportCommentReq) (*pb.ReportCommentResp, error) {
	if in.CommentId == 0 {
		return nil, fmt.Errorf("被举报评论Id不能为空")
	}
	report := &model.ReportRecord{
		UserId:     in.UserId,
		CommentId:  in.CommentId,
		TargetType: in.TargetType,
		TargetId:   in.TargetId,
		Reason:     int32(in.Reason),
		Detail:     in.Detail,
	}
	err := l.svcCtx.CommentModel.InsertReport(l.ctx, report)
	if err != nil {
		return nil, err
	}
	return &pb.ReportCommentResp{
		Success: true,
	}, nil
}
