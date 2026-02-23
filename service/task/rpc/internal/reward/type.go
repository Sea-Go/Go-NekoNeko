package reward

import (
	"strconv"
	"time"
)

const (
	StreamKey = "stream:task:reward"
	GroupName = "task-reward-group"
)

type RedisEvent struct {
	RewardID int64 //幂等ID
	UID      int64
	TaskID   int64
	Ts       int64
	AddScore int64
}

func NewEvent(_uid, _taskID string, addScore int64) *RedisEvent {
	uid, _ := strconv.ParseInt(_uid, 10, 64)
	taskID, _ := strconv.ParseInt(_taskID, 10, 64)
	return &RedisEvent{
		RewardID: uid + taskID,
		UID:      uid,
		TaskID:   taskID,
		Ts:       time.Now().Unix(),
		AddScore: addScore,
	}
}
