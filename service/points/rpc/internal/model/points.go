package model

import "gorm.io/gorm"

type PointsModel struct {
	conn *gorm.DB
}

func NewPointsModel(db *gorm.DB) *PointsModel {
	return &PointsModel{
		conn: db,
	}
}
