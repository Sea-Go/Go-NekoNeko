package model

import "time"

// Likes 对应数据库 "likes" 表
type Like struct {
	BizID      int64     `gorm:"column:biz_id;primaryKey"`
	UserID     int64     `gorm:"column:user_id;primaryKey"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
}

func (Like) TableName() string {
	return "likes"
}

// CancelLikesLog 对应数据库 "cancel_likes_log" 表
type CancelLikesLog struct {
	UserID     int64     `gorm:"column:user_id;primaryKey"`
	BizID      int64     `gorm:"column:biz_id;primaryKey"`
	UpdateTime time.Time `gorm:"column:update_time;primaryKey;autoCreateTime"`
}

func (CancelLikesLog) TableName() string {
	return "cancel_likes_log"
}

// Comment 对应数据库 "comments" 表
type Comment struct {
	ID         int64     `gorm:"column:id;primaryKey"`
	BizID      int64     `gorm:"column:biz_id;index"`
	UserID     int64     `gorm:"column:user_id;index"`
	Content    string    `gorm:"column:content"`
	ParentID   int64     `gorm:"column:parent_id;default:0"`
	RootID     int64     `gorm:"column:root_id;default:0"`
	Status     int32     `gorm:"column:status;default:1"` // 1:正常
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime"`
}

func (Comment) TableName() string {
	return "comments"
}
