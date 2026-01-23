package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CommentsModel = (*customCommentsModel)(nil)

type (
	// CommentsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCommentsModel.
	CommentsModel interface {
		commentsModel
		withSession(session sqlx.Session) CommentsModel
	}

	customCommentsModel struct {
		*defaultCommentsModel
	}
)

// NewCommentsModel returns a model for the database table.
func NewCommentsModel(conn sqlx.SqlConn) CommentsModel {
	return &customCommentsModel{
		defaultCommentsModel: newCommentsModel(conn),
	}
}

func (m *customCommentsModel) withSession(session sqlx.Session) CommentsModel {
	return NewCommentsModel(sqlx.NewSqlConnFromSession(session))
}
