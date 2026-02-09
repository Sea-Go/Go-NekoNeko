package favorite_item

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repo 收藏项数据仓库
// 职责：只负责与数据库的交互，不包含业务逻辑
type Repo struct {
	db *pgxpool.Pool
}

// NewRepo 创建仓库实例
func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

// Add 添加收藏项
// 注意：这里只做数据插入，不检查重复（由上层业务逻辑检查）
func (r *Repo) Add(ctx context.Context, item *FavoriteItem) error {
	query := `
		INSERT INTO favorite_item (user_id, folder_id, object_type, object_id, sort_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		item.UserID,
		item.FolderID,
		item.ObjectType,
		item.ObjectID,
		item.SortOrder,
		time.Now(),
		time.Now(),
	).Scan(&item.ID)
}

// Delete 软删除收藏项
func (r *Repo) Delete(ctx context.Context, userID int64, objectType string, objectID int64) error {
	query := `
		UPDATE favorite_item
		SET deleted_at = $1, updated_at = $2
		WHERE user_id = $3 AND object_type = $4 AND object_id = $5 AND deleted_at IS NULL
	`
	_, err := r.db.Exec(ctx, query, time.Now(), time.Now(), userID, objectType, objectID)
	return err
}

// GetByID 根据 ID 获取收藏项
func (r *Repo) GetByID(ctx context.Context, id int64) (*FavoriteItem, error) {
	query := `
		SELECT id, user_id, folder_id, object_type, object_id, sort_order, created_at, updated_at, deleted_at
		FROM favorite_item
		WHERE id = $1 AND deleted_at IS NULL
	`
	var item FavoriteItem
	err := r.db.QueryRow(ctx, query, id).Scan(
		&item.ID, &item.UserID, &item.FolderID, &item.ObjectType, &item.ObjectID,
		&item.SortOrder, &item.CreatedAt, &item.UpdatedAt, &item.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// ListByFolder 获取某个收藏夹中的所有项目
func (r *Repo) ListByFolder(ctx context.Context, folderID int64, offset, limit int) ([]*FavoriteItem, error) {
	query := `
		SELECT id, user_id, folder_id, object_type, object_id, sort_order, created_at, updated_at, deleted_at
		FROM favorite_item
		WHERE folder_id = $1 AND deleted_at IS NULL
		ORDER BY sort_order DESC, created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, folderID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*FavoriteItem
	for rows.Next() {
		var item FavoriteItem
		err := rows.Scan(
			&item.ID, &item.UserID, &item.FolderID, &item.ObjectType, &item.ObjectID,
			&item.SortOrder, &item.CreatedAt, &item.UpdatedAt, &item.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, rows.Err()
}

// CheckExists 检查用户是否已经收藏了某个对象
// 返回：是否存在，存在的话返回对应的收藏夹ID和项目ID
func (r *Repo) CheckExists(ctx context.Context, userID int64, objectType string, objectID int64) (bool, int64, int64, error) {
	query := `
		SELECT id, folder_id
		FROM favorite_item
		WHERE user_id = $1 AND object_type = $2 AND object_id = $3 AND deleted_at IS NULL
		LIMIT 1
	`
	var itemID, folderID int64
	err := r.db.QueryRow(ctx, query, userID, objectType, objectID).Scan(&itemID, &folderID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return false, 0, 0, nil
		}
		return false, 0, 0, err
	}
	return true, folderID, itemID, nil
}

// CountByFolder 获取收藏夹中的项目数量
func (r *Repo) CountByFolder(ctx context.Context, folderID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM favorite_item WHERE folder_id = $1 AND deleted_at IS NULL`
	var count int64
	err := r.db.QueryRow(ctx, query, folderID).Scan(&count)
	return count, err
}

// DeleteByFolder 删除某个收藏夹中的所有项目（用于删除收藏夹时）
func (r *Repo) DeleteByFolder(ctx context.Context, folderID int64) error {
	query := `
		UPDATE favorite_item
		SET deleted_at = $1, updated_at = $2
		WHERE folder_id = $3 AND deleted_at IS NULL
	`
	_, err := r.db.Exec(ctx, query, time.Now(), time.Now(), folderID)
	return err
}
