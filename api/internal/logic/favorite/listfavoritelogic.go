package favorite

import (
	"context"
	"sea-try-go/service/favorite/favorite_item"
)

// ListFavoriteLogic 列表收藏项的逻辑
type ListFavoriteLogic struct {
	itemService *favorite_item.Service
}

func NewListFavoriteLogic(itemService *favorite_item.Service) *ListFavoriteLogic {
	return &ListFavoriteLogic{itemService: itemService}
}

// Execute 执行列表查询
func (l *ListFavoriteLogic) Execute(ctx context.Context, req ListFavoriteReq, userID int64) (*ListFavoriteResp, error) {
	// 参数验证
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// 调用业务逻辑层
	items, total, err := l.itemService.ListItems(ctx, userID, req.FolderID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	var itemList []*ItemInfo
	for _, item := range items {
		itemList = append(itemList, &ItemInfo{
			ID:         item.ID,
			FolderID:   item.FolderID,
			ObjectType: item.ObjectType,
			ObjectID:   item.ObjectID,
			CreatedAt:  item.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &ListFavoriteResp{
		Items:    itemList,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// ============ 类型定义 ============

// ListFavoriteReq 列表查询请求
type ListFavoriteReq struct {
	FolderID int64 `json:"folder_id,string" form:"folder_id"`
	Page     int   `json:"page" form:"page"`
	PageSize int   `json:"page_size" form:"page_size"`
}

// ListFavoriteResp 列表响应
type ListFavoriteResp struct {
	Items    []*ItemInfo `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}
