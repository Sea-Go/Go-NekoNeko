package model

import "time"

type PointsStatus int

const (
	StatusInit       PointsStatus = 0 // 开始/初始化
	StatusProcessing PointsStatus = 1 // 正在处理
	StatusQueued     PointsStatus = 2 // 排队等候
	StatusSuccess    PointsStatus = 3 // 成功
	StatusFailed     PointsStatus = 4 // 失败
)

type Points struct {
	Id          uint64       `gorm:"primaryKey"`
	Uid         int64        `gorm:"column:uid;uniqueIndex;not null"`
	AccountId   int64        `gorm:"column:account_id;not null;uniqueIndex:idx_account_user"` // 交易ID
	UserId      int64        `gorm:"column:user_id;not null;uniqueIndex:idx_account_user"`    // 交易用户ID
	Amount      int32        `gorm:"column:amount;not null"`                                  // 交易详情（+为增加，-为减少）
	Status      PointsStatus `gorm:"column:status;default:0"`                                 // 交易状态
	WrongAnswer string       `gorm:"column:wrong_answer;type:text"`                           // 错误原因，用于统计溯源
	Tracing     string       `gorm:"column:tracing;type:text"`                                // 用于追踪
	CreatedTime time.Time    `gorm:"column:created_time;autoCreateTime"`
	UpdatedTime time.Time    `gorm:"column:updated_time;autoUpdateTime"`
}

func (Points) TableName() string {
	return "points"
}
