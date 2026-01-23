package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LikesModel = (*customLikesModel)(nil)

type (
	// LikesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLikesModel.
	LikesModel interface {
		likesModel
		withSession(session sqlx.Session) LikesModel
	}

	customLikesModel struct {
		*defaultLikesModel
	}
)

// NewLikesModel returns a model for the database table.
func NewLikesModel(conn sqlx.SqlConn) LikesModel {
	return &customLikesModel{
		defaultLikesModel: newLikesModel(conn),
	}
}

func (m *customLikesModel) withSession(session sqlx.Session) LikesModel {
	return NewLikesModel(sqlx.NewSqlConnFromSession(session))
}
