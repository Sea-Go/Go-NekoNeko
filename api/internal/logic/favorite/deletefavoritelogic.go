package favorite

import (
	"context"
	"sea-try-go/service/favorite/favorite_item"
)

// DeleteFavoriteLogic 删除收藏的逻辑
type DeleteFavoriteLogic struct {
	itemService *favorite_item.Service
}

func NewDeleteFavoriteLogic(itemService *favorite_item.Service) *DeleteFavoriteLogic {
	return &DeleteFavoriteLogic{itemService: itemService}
}

// Execute 执行删除收藏
func (l *DeleteFavoriteLogic) Execute(ctx context.Context, req DeleteFavoriteReq, userID int64) error {
	// 调用业务逻辑层
	return l.itemService.DeleteItem(ctx, userID, req.ObjectType, req.ObjectID)
}

// ============ 类型定义 ============

// DeleteFavoriteReq 删除收藏的请求
type DeleteFavoriteReq struct {
	ObjectType string `json:"object_type"`
	ObjectID   int64  `json:"object_id,string"`
}

// DeleteFavoriteResp 删除响应
type DeleteFavoriteResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
