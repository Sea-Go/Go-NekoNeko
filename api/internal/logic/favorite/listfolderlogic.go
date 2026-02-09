// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package favorite

import (
	"context"

	"sea-try-go/api/internal/svc"
	"sea-try-go/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListFolderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取收藏夹列表
func NewListFolderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFolderLogic {
	return &ListFolderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListFolderLogic) ListFolder(req *types.ListFolderReq) (resp *types.ListFolderResp, err error) {
	// todo: add your logic here and delete this line

	return
}
