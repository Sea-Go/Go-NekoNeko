package model

import "gorm.io/gorm"

type LikeRecordModel interface {
	GetTotalLikeCount(authorId int64) (int64, error)
	GetBatchLikeCount(targetType string, targetIds []string) (map[string]map[int32]int64, error)
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
