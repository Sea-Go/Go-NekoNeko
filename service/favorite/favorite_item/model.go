package favorite_item

import "time"

// FavoriteItem 收藏项（对应表 favorite_item）
// 代表用户在某个收藏夹中收藏的一个具体对象
type FavoriteItem struct {
	ID         int64      `db:"id"`
	UserID     int64      `db:"user_id"`
	FolderID   int64      `db:"folder_id"`
	ObjectType string     `db:"object_type"` // "article", "video" 等类型
	ObjectID   int64      `db:"object_id"`   // 对应的对象 ID
	SortOrder  int        `db:"sort_order"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
}

// CreateItemReq 创建收藏项的请求
type CreateItemReq struct {
	FolderID   int64  `json:"folder_id,string"`
	ObjectType string `json:"object_type"`
	ObjectID   int64  `json:"object_id,string"`
}

// DeleteItemReq 删除收藏项的请求
type DeleteItemReq struct {
	ObjectType string `json:"object_type"`
	ObjectID   int64  `json:"object_id,string"`
}

// ListItemReq 列表集合项的请求
type ListItemReq struct {
	FolderID int64 `json:"folder_id,string"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

// ItemInfo 收藏项信息（用于响应）
type ItemInfo struct {
	ID         int64     `json:"id,string"`
	FolderID   int64     `json:"folder_id,string"`
	ObjectType string    `json:"object_type"`
	ObjectID   int64     `json:"object_id,string"`
	CreatedAt  time.Time `json:"created_at"`
}
