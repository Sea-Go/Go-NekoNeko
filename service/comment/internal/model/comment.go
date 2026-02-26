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

		var existCount int64
		err := tx.Model(&CommentIndex{}).Where("id = ?", msg.CommentId).Count(&existCount).Error
		if err != nil {
			return err
		}
		if existCount > 0 {
			//发现是重复消息,直接吞掉
			return nil
		}

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
			OwnerId:    msg.OwnerId,
		}
		if msg.RootId == 0 {
			newSubject.RootCount = 1
		}
		updateCols := map[string]interface{}{
			"total_count": gorm.Expr("total_count + 1"),
		}
		if msg.RootId == 0 {
			updateCols["root_count"] = gorm.Expr("root_count + 1")
		}
		err = tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "target_type"}, {Name: "target_id"}},
			DoUpdates: clause.Assignments(updateCols),
		}).Create(&newSubject).Error
		if err != nil {
			return err
		}

		//4.更新父评论的回复数
		if msg.ParentId != 0 {
			updateCols := map[string]interface{}{
				"reply_count": gorm.Expr("reply_count + 1"),
			}
			if msg.UserId == msg.OwnerId {
				updateCols["attribute"] = gorm.Expr("attribute | ?", (1 << 1))
			}
			err := tx.Model(&CommentIndex{}).Where("id = ?", msg.ParentId).Updates(updateCols).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (m *CommentModel) ManageCommentAttribute(ctx context.Context, commentId int64, bitOffset uint, isSet bool) error {
	val := (1 << bitOffset)
	var expr clause.Expr
	if isSet {
		expr = gorm.Expr("attribute | ?", val)
	} else {
		expr = gorm.Expr("attribute & ~?", val)
	}
	return m.conn.WithContext(ctx).Model(&CommentIndex{}).Where("id = ?", commentId).Update("attribute", expr).Error
}

func (m *CommentModel) UpdateSubjectState(ctx context.Context, targetType, targetId string, state int32) error {
	return m.conn.WithContext(ctx).Model(&Subject{}).
		Where("target_type = ? AND target_id = ?", targetType, targetId).
		Update("state", state).Error
}

func (m *CommentModel) GetOwnerId(ctx context.Context, targetType, targetId string) (ownerId int, err error) {
	err = m.conn.WithContext(ctx).Model(&Subject{}).
		Where("target_type = ? AND target_id = ?", targetType, targetId).
		Select("owner_id").Scan(&ownerId).Error
	return ownerId, err
}

func (m *CommentModel) InsertReport(ctx context.Context, report *ReportRecord) error {
	return m.conn.WithContext(ctx).Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(report).Error
}
