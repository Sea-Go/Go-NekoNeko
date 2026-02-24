// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"favorite-system/internal/svc"
	"favorite-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListFoldersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListFoldersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFoldersLogic {
	return &ListFoldersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListFoldersLogic) ListFolders(req *types.ListFoldersReq) (resp *types.ListFoldersResp, err error) {
	// todo: add your logic here and delete this line

	return
}
