package folder

import "time"

// Folder 收藏夹（对应表 favorite_folder）
type Folder struct {
	ID        int64      `db:"id"`
	UserID    int64      `db:"user_id"`
	Name      string     `db:"name"`
	IsPublic  bool       `db:"is_public"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
