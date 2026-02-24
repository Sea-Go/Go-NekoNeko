// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"favorite-system/internal/svc"
	"favorite-system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveArticleFromFolderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveArticleFromFolderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveArticleFromFolderLogic {
	return &RemoveArticleFromFolderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveArticleFromFolderLogic) RemoveArticleFromFolder(req *types.RemoveArticleReq) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
