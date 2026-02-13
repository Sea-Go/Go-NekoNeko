package model

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PointsModel struct {
	conn *gorm.DB
}

func NewPointsModel(db *gorm.DB) *PointsModel {
	return &PointsModel{
		conn: db,
	}
}

func (m *PointsModel) FindOneByUid(ctx context.Context, uid int64) (*Points, error) {
	var points Points
	err := m.conn.WithContext(ctx).Where("uid = ?", uid).First(&points).Error
	if err == nil {
		return &points, nil
	}
	if err == gorm.ErrRecordNotFound {
		return nil, ErrorNotFound
	}
	return nil, err
}

func (m *PointsModel) UpdateByUid(ctx context.Context, uid int64, newPoints *Points) error {
	err := m.conn.WithContext(ctx).Model(&Points{}).Where("uid = ?", uid).Updates(newPoints).Error
	return err
}

func (m *PointsModel) Insert(ctx context.Context, points *Points) error {
	err := m.conn.WithContext(ctx).Create(points).Error
	return err
}

func (m *PointsModel) DeleteByUid(ctx context.Context, uid int64) error {
	err := m.conn.WithContext(ctx).Where("uid = ?", uid).Delete(&Points{}).Error
	return err
}

// InsertOrIgnore 幂等性插入 如果遇到(accoundId,userId)存在则忽略
func (m *PointsModel) InsertOrIgnore(ctx context.Context, points *Points) error {
	err := m.conn.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "account_id"}, {Name: "user_id"}},
		DoNothing: true,
	}).Create(points).Error
	return err
}

func (m *PointsModel) FindByAccountIdAndUserId(ctx context.Context, accountId int64, userId int64) (*Points, error) {
	var points Points
	err := m.conn.WithContext(ctx).Where("account_id = ? AND user_id = ?", accountId, userId).First(&points).Error
	return &points, err
}

func (m *PointsModel) BeginTransaction() *gorm.DB {
	return m.conn.Begin()
}

func (m *PointsModel) UpdateStatusByUid(ctx context.Context, uid int64, status PointsStatus) error {
	err := m.conn.WithContext(ctx).Model(&Points{}).Where("uid = ?", uid).Update("status", status).Error
	return err
}

func (m *PointsModel) HasProcessingByUserId(ctx context.Context, userId int64) (bool, error) {
	var count int64
	err := m.conn.WithContext(ctx).Model(&Points{}).
		Where("user_id = ? and status = ?", userId, StatusProcessing).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (m *PointsModel) HasOtherProcessingByUserId(ctx context.Context, userId int64, uid int64) (bool, error) {
	var count int64
	err := m.conn.WithContext(ctx).Model(&Points{}).
		Where("user_id = ? and status = ? and uid != ?", userId, StatusProcessing, uid).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (m *PointsModel) UpdateUserPoints(ctx context.Context, userId int64, amount int32) (bool, error) {
	result := m.conn.WithContext(ctx).Model(&User{}).
		Where("user_id = ? AND score + ? >= 0", userId, amount).
		Update("score", gorm.Expr("score + ?", amount))
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}
