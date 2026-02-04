package reward

import "time"

const (
	StreamKey = "stream:task:reward"
	GroupName = "task-reward-group"
)

type RedisEvent struct {
	RewardID int64 //幂等ID
	UID      int64
	TaskID   int64
	Ts       int64
	Score    int64
	AddScore int64
}

func NewEvent(uid int64, taskID int64, score, addScore int64) *RedisEvent {
	return &RedisEvent{
		RewardID: uid + taskID,
		UID:      uid,
		TaskID:   taskID,
		Ts:       time.Now().Unix(),
		Score:    score,
		AddScore: addScore,
	}
}
