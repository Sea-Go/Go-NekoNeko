package model

import (
	"context"

	"gorm.io/gorm"
)

type UserModel struct {
	conn *gorm.DB
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{
		conn: db,
	}
}

func (m *UserModel) FindOneByUserName(ctx context.Context, username string) (*User, error) {
	var user User
	err := m.conn.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err == nil {
		return &user, nil
	}
	if err == gorm.ErrRecordNotFound {
		return nil, ErrorNotFound
	}
	return nil, err
}

func (m *UserModel) FindOneByUid(ctx context.Context, uid int64) (*User, error) {
	var user User
	err := m.conn.WithContext(ctx).Where("uid = ?", uid).First(&user).Error
	if err == nil {
		return &user, nil
	}
	if err == gorm.ErrRecordNotFound {
		return nil, ErrorNotFound
	}
	return nil, err
}

func (m *UserModel) UpdateUserById(ctx context.Context, uid int64, newUser *User) error {
	err := m.conn.WithContext(ctx).Model(&User{}).Where("uid = ?", uid).Updates(newUser).Error
	return err
}

func (m *UserModel) Insert(ctx context.Context, user *User) error {
	err := m.conn.WithContext(ctx).Create(user).Error
	return err
}

func (m *UserModel) DeleteUserByUid(ctx context.Context, uid int64) error {
	err := m.conn.WithContext(ctx).Where("uid = ?", uid).Delete(&User{}).Error
	return err
}
