package model

import (
	"time"

	"gorm.io/gorm"
)

type LikeRecord struct {
	ID         int64  `gorm:"primaryKey;autoIncrement;comment:主键ID"`
	UserID     int64  `gorm:"type:bigint;not null;uniqueIndex:uk_user_target,priority:1;index:idx_user_list,priority:1;comment:点赞者ID"`
	TargetType string `gorm:"type:varchar(32);not null;uniqueIndex:uk_user_target,priority:2;index:idx_target_list,priority:1;comment:目标类型"`
	TargetID   string `gorm:"type:varchar(64);not null;uniqueIndex:uk_user_target,priority:3;index:idx_target_list,priority:2;comment:目标ID"`
	AuthorID   int64  `gorm:"type:bigint;not null;index:idx_author;comment:被点赞内容作者ID"`
	State      int32  `gorm:"type:smallint;not null;default:0;comment:状态(0表示无,1表示已点赞,2表示已点踩)"`

	CreatedAt time.Time      `gorm:"autoCreateTime;not null;comment:首次操作时间"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;not null;index:idx_user_list,priority:2;comment:最后操作时间"`
	DeleteAt  gorm.DeletedAt `gorm:"index;comment:GORM软删除时间"`
}

func (LikeRecord) TableName() string {
	return "like_record"
}

//uniqueIndex是唯一索引,目的是防止重名,在这里的顺序是(UserID,TargetType,TargetID),防止重复
