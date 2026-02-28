// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package folder

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"favorite-system/api/internal/svc"
	"favorite-system/api/internal/types"

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
	userId := req.UserId
	cacheKey := fmt.Sprintf("favorite:user:%d:folders", userId)

	// Try cache
	val, err := l.svcCtx.Redis.Get(l.ctx, cacheKey)
	if err == nil {
		var folders []types.FolderInfo
		if err := json.Unmarshal([]byte(val), &folders); err == nil {
			return &types.ListFoldersResp{Folders: folders}, nil
		}
	}

	// DB
	dbFolders, err := l.svcCtx.DB.ListFoldersByUser(l.ctx, userId)
	if err != nil {
		return nil, err
	}

	var folders []types.FolderInfo
	for _, f := range dbFolders {
		folders = append(folders, types.FolderInfo{
			Id:        f.ID,
			UserId:    f.UserID,
			Name:      f.Name,
			IsPublic:  f.IsPublic,
			CreatedAt: f.CreatedAt.Time.UnixMilli(),
		})
	}

	// Set cache
	if bytes, err := json.Marshal(folders); err == nil {
		_ = l.svcCtx.Redis.Set(l.ctx, cacheKey, bytes, time.Hour)
	}

	return &types.ListFoldersResp{Folders: folders}, nil
}
