package postgres

import "gorm.io/gorm"

type Favorite struct {
	gorm.Model
	UserID    uint64 `gorm:"not null;index:idx_user_article,unique"`
	ArticleID uint64 `gorm:"not null;index:idx_user_article,unique"`
}
