package model

import "time"

type Admin struct {
	Id         uint64            `gorm:"primaryKey"`
	Uid        int64             `gorm:"column:uid;uniqueIndex;not null"`
	Username   string            `gorm:"column:username;unique"`
	Password   string            `gorm:"column:password"`
	Email      string            `gorm:"column:email;unique"`
	ExtraInfo  map[string]string `gorm:"column:extra_info;serializer:json"`
	CreateTime time.Time         `gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time         `gorm:"column:update_time;autoUpdateTime"`
}

func (Admin) TableName() string {
	return "admins"
}

type User struct {
	Id         uint64            `gorm:"primaryKey"`
	Uid        int64             `gorm:"column:uid;uniqueIndex;not null"`
	Username   string            `gorm:"column:username;unique"`
	Password   string            `gorm:"column:password"`
	Email      string            `gorm:"column:email;unique"`
	Status     int64             `gorm:"column:status;default:0"`
	Score      int32             `gorm:"column:score"`
	ExtraInfo  map[string]string `gorm:"column:extra_info;serializer:json"`
	CreateTime time.Time         `gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time         `gorm:"column:update_time;autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}
