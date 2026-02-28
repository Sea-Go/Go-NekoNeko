// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package item

import (
	"context"
	"fmt"
	"strconv"

	"favorite-system/api/internal/svc"
	"favorite-system/api/internal/types"
	"favorite-system/internal/repo/db"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveFavoriteItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveFavoriteItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveFavoriteItemLogic {
	return &RemoveFavoriteItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveFavoriteItemLogic) RemoveFavoriteItem(req *types.RemoveFavoriteItemReq) (resp *types.BaseResp, err error) {
	userIdVal := l.ctx.Value("userId")
	if userIdVal == nil {
		userIdVal = l.ctx.Value("uid")
	}
	if userIdVal == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	userId, err := strconv.ParseInt(fmt.Sprintf("%v", userIdVal), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	// Check folder ownership
	folder, err := l.svcCtx.DB.GetFolderByID(l.ctx, req.FolderId)
	if err != nil {
		return nil, fmt.Errorf("folder not found")
	}
	if folder.UserID != userId {
		return nil, fmt.Errorf("permission denied")
	}

	err = l.svcCtx.DB.SoftDeleteFavoriteItem(l.ctx, db.SoftDeleteFavoriteItemParams{
		FolderID:   req.FolderId,
		UserID:     userId,
		ObjectType: req.ObjectType,
		ObjectID:   req.ObjectId,
	})
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("favorite:folder:%d:items", req.FolderId)
	_ = l.svcCtx.Redis.Del(l.ctx, cacheKey)

	return &types.BaseResp{
		Code:    0,
		Message: "success",
	}, nil
}
