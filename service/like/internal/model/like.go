package model

import "gorm.io/gorm"

type LikeRecordModel interface {
	GetTotalLikeCount(authorId int64) (int64, error)
	GetBatchLikeCount(targetType string, targetIds []string) (map[string]map[int32]int64, error)
	GetUserBatchLikeState(userId int64, targetType string, targetIds []string) (map[string]int32, error)
	GetUserLikeList(userId int64, targetType string, cursor int64, limit int) ([]UserLikeListResult, error)
}

type defaultLikeRecordModel struct {
	db *gorm.DB
}

func NewLikeRecordModel(db *gorm.DB) LikeRecordModel {
	return &defaultLikeRecordModel{db: db}
}

func (m *defaultLikeRecordModel) GetTotalLikeCount(authorId int64) (int64, error) {
	var count int64
	err := m.db.Model(&LikeRecord{}).Where("author_id = ? AND state = ?", authorId, 1).Count(&count).Error
	return count, err
}

func (m *defaultLikeRecordModel) GetBatchLikeCount(targetType string, targetIds []string) (map[string]map[int32]int64, error) {
	type Result struct {
		TargetID string
		State    int32
		Count    int64
	}
	var results []Result
	err := m.db.Model(&LikeRecord{}).
		Select("target_id,state,count(1) as count").
		Where("target_type = ? AND target_id IN (?) AND state IN (1,2)", targetType, targetIds).
		Group("target_id,state").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	resMap := make(map[string]map[int32]int64)
	for _, r := range results {
		if resMap[r.TargetID] == nil {
			resMap[r.TargetID] = make(map[int32]int64)
		}
		resMap[r.TargetID][r.State] = r.Count
	}
	return resMap, nil
}

func (m *defaultLikeRecordModel) GetUserBatchLikeState(userId int64, targetType string, targetIds []string) (map[string]int32, error) {
	type Result struct {
		TargetID string
		State    int32
	}
	var results []Result
	err := m.db.Model(&LikeRecord{}).
		Select("target_id, state").
		Where("user_id = ? AND target_type = ? AND target_id IN (?)", userId, targetType, targetIds).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	resMap := make(map[string]int32)
	for _, r := range results {
		resMap[r.TargetID] = r.State
	}
	return resMap, nil
}

type UserLikeListResult struct {
	Id         int64
	TargetId   string
	CreateTime int64
}

func (m *defaultLikeRecordModel) GetUserLikeList(userId int64, targetType string, cursor int64, limit int) ([]UserLikeListResult, error) {
	var results []UserLikeListResult
	//ID越大表示时间越新,只需要 id < cursor ,然后结合Limit就能取出最接近cursor的limit条数据
	query := m.db.Model(&LikeRecord{}).Where("user_id = ? AND target_type = ? AND state = 1", userId, targetType)
	if cursor > 0 {
		query = query.Where("id < ?", cursor)
	}
	err := query.Order("id DESC").
		Limit(limit).
		Select("id, target_id, create_time").
		Scan(&results).Error
	return results, err
}
