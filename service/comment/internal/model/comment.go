package model

import (
	"context"
	"errors"
	"fmt"
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

func (m *CommentModel) DeleteCommentTx(ctx context.Context, commentId, userId int64, targetType, targetId string) (int64, error) {
	var remainCount int64
	err := m.conn.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var comment CommentIndex
		if err := tx.Where("id = ?", commentId).First(&comment).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("评论不存在")
			}
			return err
		}
		var sub Subject
		if err := tx.Where("target_type = ? AND target_id = ?", targetType, targetId).First(&sub).Error; err != nil {
			return err
		}
		if comment.State == 2 {
			remainCount = sub.TotalCount
			return nil
		}
		if userId != comment.UserId && userId != sub.OwnerId {
			return fmt.Errorf("无权删除他人评论")
		}
		if err := tx.Model(&CommentIndex{}).Where("id = ?", commentId).Update("state", 2).Error; err != nil {
			return err
		}
		updateSub := map[string]interface{}{
			"total_count": gorm.Expr("total_count - 1"),
		}
		if comment.RootId == 0 {
			updateSub["root_count"] = gorm.Expr("root_count - 1")
		}
		if err := tx.Model(&Subject{}).Where("target_type = ? AND target_id = ?", targetType, targetId).Updates(updateSub).Error; err != nil {
			return err
		}
		tx.Where("target_type = ? AND target_id = ?", targetType, targetId).First(&sub)
		remainCount = sub.TotalCount
		if comment.ParentId != 0 {
			if err := tx.Model(&CommentIndex{}).Where("id = ?", comment.ParentId).Update("reply_count", gorm.Expr("reply_count - 1")).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return remainCount, err
}

func (m *CommentModel) GetSubjectByID(ctx context.Context, subjectId string) (Subject, error) {
	var subject Subject
	if subjectId == "" {
		return subject, fmt.Errorf("invalid subjectId empty")
	}

	err := m.conn.WithContext(ctx).
		Model(&Subject{}).
		Where("id = ?", subjectId).
		First(&subject).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Subject{}, gorm.ErrRecordNotFound
		}
		return Subject{}, fmt.Errorf("GetSubjectByID failed, subjectId=%d: %w", subjectId, err)
	}
	return subject, nil
}

func (m *CommentModel) GetReplyIDsByPage(ctx context.Context, req GetReplyIDsPageReq) ([]int64, error) {
	if req.TargetType == "" {
		return nil, fmt.Errorf("invalid TargetType: empty")
	}
	if req.TargetId == "" {
		return nil, fmt.Errorf("invalid TargetId: empty")
	}
	if req.RootId < 0 {
		return nil, fmt.Errorf("invalid RootId: %d", req.RootId)
	}
	if req.Offset < 0 {
		return nil, fmt.Errorf("invalid Offset: %d", req.Offset)
	}
	if req.Limit <= 0 {
		return nil, fmt.Errorf("invalid Limit: %d", req.Limit)
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	sort := req.Sort
	if sort == "" {
		sort = ReplySortTime
	}

	db := m.conn.WithContext(ctx).
		Model(&CommentIndex{}).
		Select("id").
		Where("target_type = ? AND target_id = ?", req.TargetType, req.TargetId).
		Where("root_id = ?", req.RootId)

	// 只查正常状态（0=正常）
	if req.OnlyNormal {
		db = db.Where("state = ?", 0)
	}

	// 排序：
	switch sort {
	case ReplySortHot:
		db = db.Order("like_count DESC").Order("id DESC")
	case ReplySortTime:
		fallthrough
	default:
		db = db.Order("created_at DESC").Order("id DESC")
	}

	db = db.Offset(req.Offset).Limit(req.Limit)

	var rows []struct {
		Id int64 `gorm:"column:id"`
	}
	if err := db.Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("GetReplyIDsByPage query failed: %w", err)
	}

	ids := make([]int64, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.Id)
	}

	return ids, nil
}

/*func (m *CommentModel) GetReplyContent(ctx context.Context, commentId int64, bitOffset uint, isSet bool) (CommentContent, error) {

}*/

func (m *CommentModel) BatchGetReplyContent(ctx context.Context, commentIDs []int64) ([]CommentContent, error) {
	if len(commentIDs) == 0 {
		return []CommentContent{}, nil
	}

	//去掉非法ID
	uniq := make(map[int64]struct{}, len(commentIDs))
	filteredIDs := make([]int64, 0, len(commentIDs))
	for _, id := range commentIDs {
		if id <= 0 {
			continue
		}
		if _, ok := uniq[id]; ok {
			continue
		}
		uniq[id] = struct{}{}
		filteredIDs = append(filteredIDs, id)
	}

	if len(filteredIDs) == 0 {
		return []CommentContent{}, nil
	}

	var rows []CommentContent
	err := m.conn.WithContext(ctx).
		Model(&CommentContent{}).
		Where("comment_id IN ?", filteredIDs).
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("BatchGetReplyContentByCommentIDs query failed: %w", err)
	}

	//便于按输入顺序重排
	contentMap := make(map[int64]CommentContent, len(rows))
	for _, row := range rows {
		contentMap[row.CommentId] = row
	}

	result := make([]CommentContent, 0, len(commentIDs))
	for _, id := range commentIDs {
		if id <= 0 {
			continue
		}
		if c, ok := contentMap[id]; ok {
			result = append(result, c)
		}
	}
	return result, nil
}
