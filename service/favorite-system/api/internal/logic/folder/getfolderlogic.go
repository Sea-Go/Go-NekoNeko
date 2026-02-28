// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package folder

import (
	"context"
	"fmt"
	"strconv"

	"favorite-system/api/internal/svc"
	"favorite-system/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFolderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFolderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFolderLogic {
	return &GetFolderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFolderLogic) GetFolder(req *types.GetFolderReq) (resp *types.GetFolderResp, err error) {
	userIdVal := l.ctx.Value("userId")
	if userIdVal == nil {
		userIdVal = l.ctx.Value("uid")
	}
	var currentUserId int64
	if userIdVal != nil {
		currentUserId, _ = strconv.ParseInt(fmt.Sprintf("%v", userIdVal), 10, 64)
	}

	folder, err := l.svcCtx.DB.GetFolderByID(l.ctx, req.FolderId)
	if err != nil {
		return nil, err
	}

	if !folder.IsPublic {
		if currentUserId == 0 || folder.UserID != currentUserId {
			return nil, fmt.Errorf("permission denied")
		}
	}

	return &types.GetFolderResp{
		Folder: types.FolderInfo{
			Id:        folder.ID,
			UserId:    folder.UserID,
			Name:      folder.Name,
			IsPublic:  folder.IsPublic,
			CreatedAt: folder.CreatedAt.Time.UnixMilli(),
		},
	}, nil
}
