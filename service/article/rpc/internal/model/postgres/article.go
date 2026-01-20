package model

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// Article GORM Model
type Article struct {
	ID            string         `gorm:"primaryKey;type:varchar(32)"`
	Title         string         `gorm:"type:varchar(255);not null"`
	Brief         string         `gorm:"type:varchar(512)"`
	Content       string         `gorm:"type:text"` // 对应 markdown_content
	CoverImageURL string         `gorm:"type:varchar(255)"`
	ManualTypeTag string         `gorm:"type:varchar(64);index"`
	SecondaryTags StringArray    `gorm:"type:jsonb"` // 使用 jsonb 存储标签数组
	AuthorID      string         `gorm:"type:varchar(32);index"`
	Status        int32          `gorm:"type:smallint;default:0"`
	ViewCount     int32          `gorm:"default:0"`
	LikeCount     int32          `gorm:"default:0"`
	CommentCount  int32          `gorm:"default:0"`
	ShareCount    int32          `gorm:"default:0"`
	ExtInfo       JSONMap        `gorm:"type:jsonb"` // 使用 jsonb 存储扩展信息
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// --- Custom Types for Postgres JSONB ---

// StringArray handles []string <-> jsonb
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, a)
}

// JSONMap handles map[string]string <-> jsonb
type JSONMap map[string]string

func (m JSONMap) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*m = nil
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, m)
}

// --- Interface & Implementation ---

type ArticleModel interface {
	Insert(ctx context.Context, article *Article) error
	FindOne(ctx context.Context, id string) (*Article, error)
	Update(ctx context.Context, article *Article) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*Article, int64, error)
}

type defaultArticleModel struct {
	db *gorm.DB
}

func NewArticleModel(db *gorm.DB) ArticleModel {
	return &defaultArticleModel{
		db: db,
	}
}

func (m *defaultArticleModel) Insert(ctx context.Context, article *Article) error {
	return m.db.WithContext(ctx).Create(article).Error
}

func (m *defaultArticleModel) FindOne(ctx context.Context, id string) (*Article, error) {
	var article Article
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (m *defaultArticleModel) Update(ctx context.Context, article *Article) error {
	return m.db.WithContext(ctx).Save(article).Error
}

func (m *defaultArticleModel) Delete(ctx context.Context, id string) error {
	return m.db.WithContext(ctx).Delete(&Article{}, "id = ?", id).Error
}

func (m *defaultArticleModel) List(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*Article, int64, error) {
	var articles []*Article
	var total int64

	query := m.db.WithContext(ctx).Model(&Article{})

	// Dynamic filtering
	if val, ok := filters["author_id"]; ok && val != "" {
		query = query.Where("author_id = ?", val)
	}
	if val, ok := filters["manual_type_tag"]; ok && val != "" {
		query = query.Where("manual_type_tag = ?", val)
	}
	if val, ok := filters["status"]; ok {
		query = query.Where("status = ?", val)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&articles).Error
	if err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}
