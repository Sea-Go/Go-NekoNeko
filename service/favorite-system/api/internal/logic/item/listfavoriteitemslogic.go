// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package item

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"favorite-system/api/internal/svc"
	"favorite-system/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListFavoriteItemsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListFavoriteItemsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFavoriteItemsLogic {
	return &ListFavoriteItemsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListFavoriteItemsLogic) ListFavoriteItems(req *types.ListFavoriteItemsReq) (resp *types.ListFavoriteItemsResp, err error) {
	userIdVal := l.ctx.Value("userId")
	if userIdVal == nil {
		userIdVal = l.ctx.Value("uid")
	}
	// userId might be nil if not logged in (if optional auth), but jwt:Auth ensures it.
	// However, we need to handle error if parsing fails.
	var currentUserId int64
	if userIdVal != nil {
		currentUserId, _ = strconv.ParseInt(fmt.Sprintf("%v", userIdVal), 10, 64)
	}

	// Check folder permission
	folder, err := l.svcCtx.DB.GetFolderByID(l.ctx, req.FolderId)
	if err != nil {
		return nil, fmt.Errorf("folder not found")
	}

	if !folder.IsPublic {
		if currentUserId == 0 || folder.UserID != currentUserId {
			return nil, fmt.Errorf("permission denied")
		}
	}

	// Pagination
	limit := req.PageSize
	offset := (req.Page - 1) * req.PageSize

	// Cache (Try to get ALL items and slice in memory)
	cacheKey := fmt.Sprintf("favorite:folder:%d:items", req.FolderId)
	val, err := l.svcCtx.Redis.Get(l.ctx, cacheKey)
	if err == nil {
		var allItems []types.FavoriteItemInfo
		if err := json.Unmarshal([]byte(val), &allItems); err == nil {
			total := int64(len(allItems))
			start := int64(offset)
			end := start + int64(limit)
			if start >= total {
				return &types.ListFavoriteItemsResp{
					Items:    []types.FavoriteItemInfo{},
					Total:    total,
					Page:     req.Page,
					PageSize: req.PageSize,
				}, nil
			}
			if end > total {
				end = total
			}
			return &types.ListFavoriteItemsResp{
				Items:    allItems[start:end],
				Total:    total,
				Page:     req.Page,
				PageSize: req.PageSize,
			}, nil
		}
	}

	// DB Fallback
	allItems, err := l.svcCtx.DB.ListAllFavoriteItems(l.ctx, req.FolderId)
	if err != nil {
		return nil, err
	}

	var allRespItems []types.FavoriteItemInfo
	for _, i := range allItems {
		allRespItems = append(allRespItems, types.FavoriteItemInfo{
			Id:         i.ID,
			FolderId:   i.FolderID,
			UserId:     i.UserID,
			ObjectType: i.ObjectType,
			ObjectId:   i.ObjectID,
			Title:      i.Title,
			CreatedAt:  i.CreatedAt.Time.UnixMilli(),
		})
	}

	// Set Cache
	if bytes, err := json.Marshal(allRespItems); err == nil {
		_ = l.svcCtx.Redis.Set(l.ctx, cacheKey, bytes, time.Hour)
	}

	// Slice for response
	total := int64(len(allRespItems))
	start := int64(offset)
	end := start + int64(limit)

	if start >= total {
		return &types.ListFavoriteItemsResp{
			Items:    []types.FavoriteItemInfo{},
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		}, nil
	}
	if end > total {
		end = total
	}

	return &types.ListFavoriteItemsResp{
		Items:    allRespItems[start:end],
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
