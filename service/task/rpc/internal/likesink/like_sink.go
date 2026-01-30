package likesink

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserLikeCount struct {
	UserId    string    `gorm:"primary_key;column:user_id"`
	LikeCount int64     `gorm:"column:like_count"`
	UpdateAt  time.Time `gorm:"column:update_at"`
}

func (UserLikeCount) TableName() string {
	return "user_likes_count"
}

type LikeSinkConsumer struct {
	rdb *redis.Client
	gdb *gorm.DB

	mu    sync.Mutex
	delta map[string]int64

	flushEvery   time.Duration
	redisPipeMax int
	pgBatchMax   int

	flushCh chan struct{}
}

type row struct {
	uid   string
	delta int64
}

func NewSinkConsumer(rdb *redis.Client, gdb *gorm.DB) *LikeSinkConsumer {
	return &LikeSinkConsumer{
		rdb:          rdb,
		gdb:          gdb,
		delta:        make(map[string]int64, 1<<16),
		flushEvery:   1 * time.Second,
		redisPipeMax: 5000,
		pgBatchMax:   2000,
		flushCh:      make(chan struct{}, 1),
	}
}

func (c *LikeSinkConsumer) Start(ctx context.Context) {
	go c.loop(ctx)
}

func (c *LikeSinkConsumer) Consume(ctx context.Context, key string, val string) error {
	userID := key
	if userID == "" {
		return nil
	}

	d := int64(1)
	if val != "" {
		if n, err := strconv.ParseInt(val, 10, 64); err == nil {
			d = n
		}
	}
	if d == 0 {
		return nil
	}

	c.mu.Lock()
	c.delta[userID] += d
	c.mu.Unlock()
	return nil
}

func (c *LikeSinkConsumer) loop(ctx context.Context) {

	ticker := time.NewTicker(c.flushEvery)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			_ = c.flushOnce(context.Background())
			return
		case <-ticker.C:
			_ = c.flushOnce(ctx)
		case <-c.flushCh:
			_ = c.flushOnce(ctx)
		}
	}
}

func (c *LikeSinkConsumer) flushOnce(ctx context.Context) error {
	batch := c.swap()
	if len(batch) == 0 {
		return nil
	}
	if err := c.flushRedis(ctx, batch); err != nil {
		return err
	}
	if err := c.flushPostgres(ctx, batch); err != nil {
		return err
	}
	return nil
}

func (c *LikeSinkConsumer) swap() map[string]int64 {

	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.delta) == 0 {
		return nil
	}
	b := c.delta
	c.delta = make(map[string]int64, 1<<16)
	return b
}

func (c *LikeSinkConsumer) flushRedis(ctx context.Context, batch map[string]int64) error {
	pipe := c.rdb.Pipeline()
	n := 0

	exec := func() error {
		if n == 0 {
			return nil
		}
		_, err := pipe.Exec(ctx)
		pipe = c.rdb.Pipeline()
		n = 0
		return err
	}

	for uid, d := range batch {
		pipe.IncrBy(ctx, "like:total:"+uid, d)
		n++
		if n > c.redisPipeMax {
			if err := exec(); err != nil {
				return err
			}
		}
	}
	return exec()
}

func (c *LikeSinkConsumer) flushPostgres(ctx context.Context, batch map[string]int64) error {

	rows := make([]row, 0, len(batch))
	for uid, d := range batch {
		rows = append(rows, row{uid, d})
	}

	for i := 0; i < len(rows); i++ {
		end := i + c.pgBatchMax
		if end > len(rows) {
			end = len(rows)
		}
		if err := c.upsertChunk(rows[i:end]); err != nil {
			return err
		}
	}
	return nil
}

func (c *LikeSinkConsumer) upsertChunk(rows []row) error {
	records := make([]UserLikeCount, 0, len(rows))
	for _, row := range rows {
		records = append(records, UserLikeCount{
			UserId:    row.uid,
			LikeCount: row.delta,
		})
	}

	return c.gdb.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.Assignments(map[string]any{
				"like_count": gorm.Expr("user_like_count.like_count + EXCLUDE.like_count"),
				"update_at":  gorm.Expr("now()"),
			}),
		}).Create(&records).Error
}
