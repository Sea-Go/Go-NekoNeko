package folder

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

// Create 创建收藏夹
func (r *Repo) Create(ctx context.Context, f *Folder) error {
	sql := `
		INSERT INTO favorite_folder (user_id, name, is_public)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		ctx,
		sql,
		f.UserID,
		f.Name,
		f.IsPublic,
	).Scan(
		&f.ID,
		&f.CreatedAt,
		&f.UpdatedAt,
	)
}

// ExistsByUserAndName 判断是否重名
func (r *Repo) ExistsByUserAndName(
	ctx context.Context,
	userID int64,
	name string,
) (bool, error) {

	sql := `
		SELECT EXISTS (
			SELECT 1
			FROM favorite_folder
			WHERE user_id = $1
			  AND name = $2
			  AND deleted_at IS NULL
		)
	`

	var exists bool
	err := r.db.QueryRow(ctx, sql, userID, name).Scan(&exists)
	return exists, err
}

// SoftDelete 软删除
func (r *Repo) SoftDelete(ctx context.Context, id int64) error {
	sql := `
		UPDATE favorite_folder
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := r.db.Exec(ctx, sql, id)
	return err
}

// ListByUser 返回用户的收藏夹列表
func (r *Repo) ListByUser(ctx context.Context, userID int64) ([]*Folder, error) {
	sql := `
		SELECT id, user_id, name, is_public, created_at, updated_at, deleted_at
		FROM favorite_folder
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY id DESC
	`

	rows, err := r.db.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Folder
	for rows.Next() {
		f := &Folder{}
		if err := rows.Scan(
			&f.ID,
			&f.UserID,
			&f.Name,
			&f.IsPublic,
			&f.CreatedAt,
			&f.UpdatedAt,
			&f.DeletedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, f)
	}
	return list, rows.Err()
}
