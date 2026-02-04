package reward

import "time"

const (
	StreamKey = "stream:task:reward"
	GroupName = "task-reward-group"
)

type RedisEvent struct {
	RewardID string
	UID      string
	TaskID   string
	Ts       int64
	Score    int64
}

func NewEvent(uid, taskID string, score int64) *RedisEvent {
	return &RedisEvent{
		RewardID: uid + ":" + taskID,
		UID:      uid,
		TaskID:   taskID,
		Ts:       time.Now().Unix(),
		Score:    score,
	}
}
