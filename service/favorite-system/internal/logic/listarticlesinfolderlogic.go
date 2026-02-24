// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"favorite-system/internal/svc"
	"favorite-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListArticlesInFolderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListArticlesInFolderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListArticlesInFolderLogic {
	return &ListArticlesInFolderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListArticlesInFolderLogic) ListArticlesInFolder(req *types.ListArticlesReq) (resp *types.ListArticlesResp, err error) {
	// todo: add your logic here and delete this line

	return
}
