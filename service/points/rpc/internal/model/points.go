package model

import (
	"context"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PointsModel struct {
	conn *gorm.DB
}

func NewPointsModel(db *gorm.DB) *PointsModel {
	return &PointsModel{
		conn: db,
	}
}

// 用户
func (m *PointsModel) UpdateScore(ctx context.Context, uid int64, score decimal.Decimal) error {
	err := m.conn.WithContext(ctx).Model(&User{}).Where("uid = ?", uid).Update("score", gorm.Expr("score + ?", score)).Error
	return err
}
func (m *PointsModel) FindUserById(ctx context.Context, uid int64) (user User, err error) {
	err = m.conn.WithContext(ctx).Model(&User{}).Where("uid = ?", uid).First(&user).Error
	return user, err
}

// 流水
func (m *PointsModel) CreateTransaction(ctx context.Context, transaction *Transaction) error {
	err := m.conn.WithContext(ctx).Create(&transaction).Error
	return err
}

func (m *PointsModel) UpdateTransactionStatus(ctx context.Context, accountId string, userId int64, status int64, wrong string) error {
	err := m.conn.WithContext(ctx).Model(&Transaction{}).Where("account_id = ? AND user_id = ?", accountId, userId).Updates(
		map[string]interface{}{
			"status":       status,
			"wrong_answer": wrong,
		}).Error

	return err
}
func (m *PointsModel) CountRunningByUserId(ctx context.Context, userId int64) (int64, error) {
	var res int64
	err := m.conn.WithContext(ctx).Model(&Transaction{}).Where("user_id = ? AND status = 1", userId).Count(&res).Error
	return res, err
}
func (m *PointsModel) FindByAccountIdAndUserId(ctx context.Context, accountId string, userId int64) (*Transaction, error) {
	var tx Transaction
	err := m.conn.WithContext(ctx).Model(&Transaction{}).Where("account_id = ? AND user_id = ?", accountId, userId).First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}
