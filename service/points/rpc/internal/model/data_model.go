package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	Id         uint64            `gorm:"primaryKey"`
	Uid        int64             `gorm:"column:uid;uniqueIndex;not null"`
	Username   string            `gorm:"column:username;unique"`
	Password   string            `gorm:"column:password"`
	Email      string            `gorm:"column:email;unique"`
	Status     int64             `gorm:"column:status;default:0"`
	Score      decimal.Decimal   `gorm:"column:score"`
	ExtraInfo  map[string]string `gorm:"column:extra_info;serializer:json"`
	CreateTime time.Time         `gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time         `gorm:"column:update_time;autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

type Transaction struct {
	AccountId   string          `gorm:"primaryKey;column:account_id;"`
	UserId      int64           `gorm:"primaryKey;column:user_id;"`
	Amount      decimal.Decimal `gorm:"column:amount;not null"`
	Status      int64           `gorm:"column:status;default:0"` //0:排队中,1:处理中,2:成功,-1:失败
	WrongAnswer string          `gorm:"column:wrong_answer;"`
	Tracing     string          `gorm:"column:tracing;"`
	RetryCount  int64           `gorm:"column:retry_count;default:0"`
	CreateTime  time.Time       `gorm:"column:create_time;autoCreateTime"`
	UpdateTime  time.Time       `gorm:"column:update_time;autoUpdateTime"`
}

func (Transaction) TableName() string {
	return "transactions"
}
