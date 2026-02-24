package reward

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Product struct {
	rdb *redis.Client
}

func NewProduct(rdb *redis.Client) *Product {
	return &Product{
		rdb: rdb,
	}
}

func (p *Product) Enqueue(ctx context.Context, ev *RedisEvent) error {
	_, err := p.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: StreamKey,
		Values: map[string]any{
			"reward_id": ev.RewardID,
			"uid":       ev.UID,
			"task_id":   ev.TaskID,
			"ts":        ev.Ts,
			"score":     ev.AddScore,
		},
	}).Result()
	return err
}
