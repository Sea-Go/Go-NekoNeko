package favorite

import (
	"context"
	"sea-try-go/service/favorite/favorite_item"
)

// CreateFavoriteLogic 添加收藏的逻辑
type CreateFavoriteLogic struct {
	itemService *favorite_item.Service
}

func NewCreateFavoriteLogic(itemService *favorite_item.Service) *CreateFavoriteLogic {
	return &CreateFavoriteLogic{itemService: itemService}
}

// Execute 执行添加收藏
func (l *CreateFavoriteLogic) Execute(ctx context.Context, req CreateFavoriteReq, userID int64) (*ItemInfo, error) {
	// 调用业务逻辑层
	item, err := l.itemService.CreateItem(ctx, userID, req.FolderID, req.ObjectType, req.ObjectID)
	if err != nil {
		return nil, err
	}
	return &ItemInfo{
		ID:         item.ID,
		FolderID:   item.FolderID,
		ObjectType: item.ObjectType,
		ObjectID:   item.ObjectID,
		CreatedAt:  item.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// ============ 类型定义 ============

// CreateFavoriteReq 添加收藏的请求
type CreateFavoriteReq struct {
	FolderID   int64  `json:"folder_id,string"`
	ObjectType string `json:"object_type"`
	ObjectID   int64  `json:"object_id,string"`
}

// ItemInfo 收藏项信息
type ItemInfo struct {
	ID         int64  `json:"id,string"`
	FolderID   int64  `json:"folder_id,string"`
	ObjectType string `json:"object_type"`
	ObjectID   int64  `json:"object_id,string"`
	CreatedAt  string `json:"created_at"`
}
