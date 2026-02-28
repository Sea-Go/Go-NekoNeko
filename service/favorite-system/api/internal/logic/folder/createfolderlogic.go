// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package folder

import (
	"context"
	"fmt"
	"strconv"

	"favorite-system/api/internal/svc"
	"favorite-system/api/internal/types"
	"favorite-system/internal/repo/db"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateFolderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateFolderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFolderLogic {
	return &CreateFolderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateFolderLogic) CreateFolder(req *types.CreateFolderReq) (resp *types.CreateFolderResp, err error) {
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

	folder, err := l.svcCtx.DB.CreateFolder(l.ctx, db.CreateFolderParams{
		UserID:   userId,
		Name:     req.Name,
		IsPublic: req.IsPublic,
	})
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("favorite:user:%d:folders", userId)
	_ = l.svcCtx.Redis.Del(l.ctx, cacheKey)

	return &types.CreateFolderResp{
		Id: folder.ID,
	}, nil
}
