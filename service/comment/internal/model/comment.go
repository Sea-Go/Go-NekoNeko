package model

import (
	"context"
	kqtypes "sea-try-go/service/comment/common/types"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommentModel struct {
	conn *gorm.DB
}

func NewCommentModel(db *gorm.DB) *CommentModel {
	return &CommentModel{
		conn: db,
	}
}

func (m *CommentModel) InsertCommentTx(ctx context.Context, msg kqtypes.CommentKafkaMsg, status int) error {
	//事务
	return m.conn.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//时间校准
		createTime := time.Unix(msg.CreateTime, 0)

		//1.存入正文表
		content := &CommentContent{
			CommentId: msg.CommentId,
			Content:   msg.Content,
			Meta:      msg.Meta,
			CreatedAt: createTime,
		}
		if err := tx.Create(content).Error; err != nil {
			return err
		}

		//2.存入索引表
		index := &CommentIndex{
			Id:         msg.CommentId,
			TargetType: msg.TargetType,
			TargetId:   msg.TargetId,
			UserId:     msg.UserId,
			RootId:     msg.RootId,
			ParentId:   msg.ParentId,
			State:      int32(status),
			Attribute:  msg.Attribute,
			CreatedAt:  createTime,
		}
		if err := tx.Create(index).Error; err != nil {
			return err
		}

		//3.更新Subject表
		newSubject := &Subject{
			TargetType: msg.TargetType,
			TargetId:   msg.TargetId,
			TotalCount: 1,
			RootCount:  0,
			State:      0,
			Attribute:  0,
		}
		if msg.RootId == 0 {
			newSubject.RootCount = 1
		}
		updateCols := map[string]interface{}{
			"total_count": gorm.Expr("total_count + 1"),
		}
		if msg.RootId == 0 {
			updateCols["root_count"] = gorm.Expr("subject.root_count + 1")
		}
		err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "target_type"}, {Name: "subject.target_id"}},
			DoUpdates: clause.Assignments(updateCols),
		}).Create(&newSubject).Error
		if err != nil {
			return err
		}

		//4.更新父评论的回复数
		if msg.ParentId != 0 {
			err := tx.Model(&CommentIndex{}).Where("id = ?", msg.ParentId).Update("reply_count", gorm.Expr("reply_count + 1")).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}
