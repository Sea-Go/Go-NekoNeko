package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sea-try-go/service/comment/internal/model"
	"time"

	"github.com/redis/go-redis/v9"
)

const defaultSubjectTTL = 5 * time.Minute

func (c *CommentCache) GetSubjectWithCache(ctx context.Context, subjectID string, conn *model.CommentModel) (model.Subject, error) {
	if c == nil || c.rdb == nil {
		return model.Subject{}, fmt.Errorf("comment cache is nil")
	}
	if conn == nil {
		return model.Subject{}, fmt.Errorf("comment model conn is nil")
	}
	if subjectID == "" {
		return model.Subject{}, fmt.Errorf("invalid subjectID: empty")
	}

	if cached, err := c.GetSubjectCache(ctx, subjectID); err == nil && cached != nil {
		return *cached, nil
	} else if err != nil {
		// Redis异常记日志
	}

	sfKey := fmt.Sprintf("subject:%s", subjectID)

	v, err, _ := c.sf.Do(sfKey, func() (interface{}, error) {
		if cached, err := c.GetSubjectCache(ctx, subjectID); err == nil && cached != nil {
			return *cached, nil
		}

		subject, dbErr := conn.GetSubjectByID(ctx, subjectID)
		if dbErr != nil {
			return model.Subject{}, dbErr
		}

		_ = c.SetSubjectCache(ctx, subjectID, &subject, 5*time.Minute)

		return subject, nil
	})
	if err != nil {
		return model.Subject{}, err
	}

	subject, ok := v.(model.Subject)
	if !ok {
		return model.Subject{}, fmt.Errorf("singleflight result type assert failed")
	}

	return subject, nil
}

func (c *CommentCache) GetSubjectCache(ctx context.Context, subjectID string) (*model.Subject, error) {
	if c == nil || c.rdb == nil {
		return nil, fmt.Errorf("comment cache is nil")
	}
	if subjectID == "" {
		return nil, fmt.Errorf("invalid subjectID: empty")
	}

	val, err := c.rdb.Get(ctx, SubjectKey(subjectID)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var s model.Subject
	if err := json.Unmarshal([]byte(val), &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (c *CommentCache) SetSubjectCache(ctx context.Context, subjectID string, subject *model.Subject, ttl time.Duration) error {
	if c == nil || c.rdb == nil {
		return fmt.Errorf("comment cache is nil")
	}
	if subjectID == "" {
		return fmt.Errorf("invalid subjectID: empty")
	}
	if subject == nil {
		return fmt.Errorf("subject is nil")
	}
	if ttl <= 0 {
		ttl = defaultSubjectTTL
	}

	key := SubjectKey(subjectID)

	b, err := json.Marshal(subject)
	if err != nil {
		return fmt.Errorf("marshal subject cache failed, subjectID=%s: %w", subjectID, err)
	}

	if err := c.rdb.Set(ctx, key, b, ttl).Err(); err != nil {
		return fmt.Errorf("redis set subject cache failed, key=%s: %w", key, err)
	}

	return nil
}
