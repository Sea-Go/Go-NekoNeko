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
	UserID    string    `gorm:"primary_key;column:user_id"`
	LikeCount int64     `gorm:"column:like_count"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (UserLikeCount) TableName() string {
	return "user_like_count"
}

const (
	taskID = ""
	target = ""
	total  = 5
)

type UserTaskProgress struct {
	UserID   string  `gorm:"column:user_id"`
	TaskID   string  `gorm:"primary_key;column:task_id"`
	Status   string  `gorm:"column:status"`
	Progress float64 `gorm:"column:progress"`
	Target   string  `gorm:"column:target"`
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
	if err := c.lazyInitTaskIfNeeded(ctx, batch); err != nil {
		return err
	}
	if err := c.flushRedis(ctx, batch); err != nil {
		return err
	}
	if err := c.flushPostgres(ctx, batch); err != nil {
		return err
	}
	return nil
}

func (c *LikeSinkConsumer) lazyInitTaskIfNeeded(ctx context.Context, batch map[string]int64) error {

	pipe := c.rdb.Pipeline()
	type item struct {
		uid string
		cmd *redis.BoolCmd
	}
	size := 0
	exec := func() error {
		if size == 0 {
			return nil
		}
		_, err := pipe.Exec(ctx)
		pipe = c.rdb.Pipeline()
		size = 0
		return err
	}

	items := make([]item, 0, len(batch))
	for uid := range batch {
		k := "task:init:" + uid + ":" + taskID
		size++
		items = append(items, item{uid: uid, cmd: pipe.SetNX(ctx, k, "1", 90*24*time.Hour)})
		if size >= c.redisPipeMax {
			if err := exec(); err != nil {
				return err
			}
		}
	}
	if err := exec(); err != nil {
		return err
	}

	newUsers := make([]string, 0)
	for _, it := range items {
		ok, err := it.cmd.Result()
		if err == nil && ok {
			newUsers = append(newUsers, it.uid)
		}
	}
	if len(newUsers) == 0 {
		return nil
	}

	pipe2 := c.rdb.Pipeline()
	now := time.Now().Unix()
	for _, uid := range newUsers {
		pk := "task:progress:" + uid + ":" + taskID
		pipe2.HSet(ctx, pk,
			"status", "doing",
			"process", 0,
			"target", "target",
			"createAt", now,
			"updateAt", now)
	}

	if _, err := pipe2.Exec(ctx); err != nil {
		return err
	}

	records := make([]UserTaskProgress, 0, len(newUsers))
	for _, uid := range newUsers {
		records = append(records, UserTaskProgress{
			UserID:   uid,
			TaskID:   taskID,
			Status:   "doing",
			Progress: 0,
			Target:   target,
		})
	}
	return c.gdb.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "task_id"}},
		DoNothing: true,
	}).Create(&records).Error
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

	lens := min(len(batch), c.redisPipeMax+1)
	cmds := make(map[string]*redis.IntCmd, lens)

	exec := func() error {
		if n == 0 {
			return nil
		}
		_, err := pipe.Exec(ctx)
		if err != nil {
			pipe = c.rdb.Pipeline()
			n = 0
			cmds = make(map[string]*redis.IntCmd, lens)
			return err
		}
		for uid, cmd := range cmds {
			nowTotal, e := cmd.Result()
			if e != nil {
				continue
			}
			if nowTotal > 5 {
				_ = c.completeLikeGT5(ctx, uid, nowTotal)
			}
		}

		pipe = c.rdb.Pipeline()
		n = 0
		cmds = make(map[string]*redis.IntCmd, lens)
		return nil
	}

	for uid, d := range batch {
		cmds[uid] = pipe.IncrBy(ctx, "like:total:"+uid, d)
		n++
		if n > c.redisPipeMax {
			if err := exec(); err != nil {
				return err
			}
		}
	}
	return exec()
}

func (c *LikeSinkConsumer) completeLikeGT5(ctx context.Context, uid string, total int64) error {
	doneKey := "task:done:" + uid + ":" + taskID
	ok, err := c.rdb.SetNX(ctx, doneKey, "1", 90*24*time.Hour).Result()
	if err != nil || !ok {
		return err
	}

	now := time.Now().Unix()
	pk := "task:progress:" + uid + ":" + taskID
	_, err = c.rdb.HSet(ctx, pk,
		"status", "done",
		"doneAt", now,
		"progress", total, // 你也可以写成 5 或 likeTotal，看你语义
		"updatedAt", now,
	).Result()
	return err
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
			UserID:    row.uid,
			LikeCount: row.delta,
		})
	}

	return c.gdb.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.Assignments(map[string]any{
				"like_count": gorm.Expr("user_like_count.like_count + EXCLUDED.like_count"),
				"updated_at": gorm.Expr("now()"),
			}),
		}).Create(&records).Error
}
