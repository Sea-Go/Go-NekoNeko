package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ArticleLikesModel = (*customArticleLikesModel)(nil)

/*
CREATE TABLE "public"."article_likes" (
  "id" bigserial NOT NULL,            -- 自增主键
  "article_id" varchar(64) NOT NULL,  -- 文章ID
  "user_id" varchar(64) NOT NULL,     -- 用户ID
  "create_time" timestamp DEFAULT CURRENT_TIMESTAMP, -- 点赞时间
-- 主键
  CONSTRAINT "article_likes_pkey" PRIMARY KEY ("id"),
-- 联合唯一约束防止用户重复点赞
-- 加速 "WHERE article_id=?" 和 "WHERE article_id=? AND user_id=?"
  CONSTRAINT "uk_article_user" UNIQUE ("article_id", "user_id")
);

-- 查询这个用户点赞过哪些文章
CREATE INDEX "idx_article_likes_user_id" ON "public"."article_likes" USING btree (
  "user_id"
);
*/

type (
	// ArticleLikesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customArticleLikesModel.
	ArticleLikesModel interface {
		articleLikesModel
		withSession(session sqlx.Session) ArticleLikesModel
	}

	customArticleLikesModel struct {
		*defaultArticleLikesModel
	}
)

// NewArticleLikesModel returns a model for the database table.
func NewArticleLikesModel(conn sqlx.SqlConn) ArticleLikesModel {
	return &customArticleLikesModel{
		defaultArticleLikesModel: newArticleLikesModel(conn),
	}
}

func (m *customArticleLikesModel) withSession(session sqlx.Session) ArticleLikesModel {
	return NewArticleLikesModel(sqlx.NewSqlConnFromSession(session))
}
