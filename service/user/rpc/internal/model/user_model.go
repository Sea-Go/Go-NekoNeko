package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
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

// UserModel 用户数据模型
type UserModel struct {
	db *gorm.DB
}

// NewUserModel 创建新的用户模型
func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{db: db}
}

// Create 创建用户
func (m *UserModel) Create(user *User) error {
	return m.db.Create(user).Error
}

// GetByUid 获取用户信息
func (m *UserModel) GetByUid(uid int64) (*User, error) {
	var user User
	err := m.db.Where("uid = ?", uid).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

// GetByUsername 根据用户名查询用户
func (m *UserModel) GetByUsername(username string) (*User, error) {
	var user User
	err := m.db.Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

// Update 更新用户信息
func (m *UserModel) Update(user *User) error {
	return m.db.Model(user).Updates(user).Error
}

// Delete 删除用户
func (m *UserModel) Delete(uid int64) error {
	return m.db.Where("uid = ?", uid).Delete(&User{}).Error
}
