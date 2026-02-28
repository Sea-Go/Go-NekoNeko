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

type CheckFavoriteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckFavoriteLogic {
	return &CheckFavoriteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckFavoriteLogic) CheckFavorite(req *types.CheckFavoriteReq) (resp *types.CheckFavoriteResp, err error) {
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

	item, err := l.svcCtx.DB.GetFavoriteItem(l.ctx, db.GetFavoriteItemParams{
		UserID:     userId,
		ObjectType: req.ObjectType,
		ObjectID:   req.ObjectId,
	})
	
	if err != nil {
		// Not found or db error
		// sqlc returns error if not found? Yes, sql.ErrNoRows usually.
		// We should return false instead of error if not found.
		return &types.CheckFavoriteResp{
			IsFavorited: false,
		}, nil
	}

	return &types.CheckFavoriteResp{
		IsFavorited: true,
		FolderId:    item.FolderID,
	}, nil
}
