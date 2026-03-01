// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package comment

import (
	"context"

	"sea-try-go/service/comment/api/internal/svc"
	"sea-try-go/service/comment/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReportCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReportCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportCommentLogic {
	return &ReportCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReportCommentLogic) ReportComment(req *types.ReportCommentReq) (resp *types.ReportCommentResp, err error) {
	// todo: add your logic here and delete this line

	return
}
