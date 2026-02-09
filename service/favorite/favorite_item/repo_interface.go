package favorite_item

import (
	"context"
)

// RepoInterface 提供 Repo 所需的方法签名，便于测试替换实现
type RepoInterface interface {
	Add(ctx context.Context, item *FavoriteItem) error
	Delete(ctx context.Context, userID int64, objectType string, objectID int64) error
	ListByFolder(ctx context.Context, folderID int64, offset, limit int) ([]*FavoriteItem, error)
	CheckExists(ctx context.Context, userID int64, objectType string, objectID int64) (bool, int64, int64, error)
	CountByFolder(ctx context.Context, folderID int64) (int64, error)
}
