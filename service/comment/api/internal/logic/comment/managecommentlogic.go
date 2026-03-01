// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package comment

import (
	"context"

	"sea-try-go/service/comment/api/internal/svc"
	"sea-try-go/service/comment/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ManageCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewManageCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ManageCommentLogic {
	return &ManageCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ManageCommentLogic) ManageComment(req *types.ManageCommentReq) (resp *types.ManageCommentResp, err error) {
	// todo: add your logic here and delete this line

	return
}
