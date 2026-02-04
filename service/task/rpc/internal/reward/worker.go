package reward

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type PointsClient interface {
	// rpc调用
	AddPoints(ctx context.Context, rewardID, uid, reason string, delta int64) error
}

type Worker struct {
	rdb          *redis.Client
	points       PointsClient
	consumerName string

	batch    int64
	blockDur time.Duration

	failSleep time.Duration
}

func NewWorker(rdb *redis.Client, points PointsClient, consumerName string) *Worker {
	return &Worker{
		rdb:          rdb,
		points:       points,
		consumerName: consumerName,
		batch:        64,
		blockDur:     5 * time.Second,
		failSleep:    300 * time.Millisecond,
	}
}

func (w *Worker) EnsureGroup(ctx context.Context) error {
	err := w.rdb.XGroupCreateMkStream(ctx, StreamKey, GroupName, "0").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		return err
	}
	return nil
}

func (w *Worker) Run(ctx context.Context) error {
	if err := w.EnsureGroup(ctx); err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		res, err := w.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    GroupName,
			Consumer: w.consumerName,
			Streams:  []string{StreamKey, ">"},
			Count:    w.batch,
			Block:    w.blockDur,
		}).Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			time.Sleep(w.failSleep)
			continue
		}
		for _, s := range res {
			for _, msg := range s.Messages {
				ev := parseRedisEvent(msg.Values)

				delta := ev.Score
				reason := "task_done:" + ev.TaskID

				if err := w.points.AddPoints(ctx, ev.RewardID, ev.UID, reason, delta); err != nil {
					//留给后边重试
					continue
				}
				_ = w.rdb.XAck(ctx, StreamKey, GroupName, msg.ID).Err()
				_ = w.rdb.XDel(ctx, StreamKey, msg.ID).Err()
			}
		}
	}
}

func parseRedisEvent(values map[string]any) RedisEvent {
	getStr := func(k string) string {
		if v, ok := values[k]; ok {
			switch x := v.(type) {
			case string:
				return x
			case []byte:
				return string(x)
			default:
				return ""
			}
		}
		return ""
	}
	getI64 := func(k string) int64 {
		if v, ok := values[k]; ok {
			switch x := v.(type) {
			case int64:
				return x
			case int:
				return int64(x)
			case string:
				n, _ := strconv.ParseInt(x, 10, 64)
				return n
			case []byte:
				n, _ := strconv.ParseInt(string(x), 10, 64)
				return n
			default:
				return 0
			}
		}
		return 0
	}

	return RedisEvent{
		RewardID: getStr("reward_id"),
		UID:      getStr("uid"),
		TaskID:   getStr("task_id"),
		Ts:       getI64("ts"),
		Score:    getI64("score"),
	}
}
