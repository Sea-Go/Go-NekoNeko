// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package favorite

import (
	"context"

	"sea-try-go/api/internal/svc"
	"sea-try-go/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateFolderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建收藏夹
func NewCreateFolderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFolderLogic {
	return &CreateFolderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateFolderLogic) CreateFolder(req *types.CreateFolderReq) (resp *types.CreateFolderResp, err error) {
	// todo: add your logic here and delete this line

	return
}
