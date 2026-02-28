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

type DeleteFolderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteFolderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFolderLogic {
	return &DeleteFolderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteFolderLogic) DeleteFolder(req *types.DeleteFolderReq) (resp *types.BaseResp, err error) {
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

	folder, err := l.svcCtx.DB.GetFolderByID(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if folder.UserID != userId {
		return nil, fmt.Errorf("permission denied")
	}

	err = l.svcCtx.DB.SoftDeleteFolder(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf("favorite:user:%d:folders", userId)
	_ = l.svcCtx.Redis.Del(l.ctx, cacheKey)

	return &types.BaseResp{
		Code:    0,
		Message: "success",
	}, nil
}
