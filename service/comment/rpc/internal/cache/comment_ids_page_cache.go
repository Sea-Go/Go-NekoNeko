package cache

import (
	"context"
	"encoding/json"
	"fmt"
	model2 "sea-try-go/service/comment/rpc/internal/model"
	"time"

	"github.com/redis/go-redis/v9"
)

const defaultReplyIDsPageTTL = 3 * time.Minute

func (c *CommentCache) GetReplyIDsPageCache(ctx context.Context, req model2.GetReplyIDsPageReq, conn *model2.CommentModel) ([]int64, error) {
	if c == nil || c.rdb == nil {
		return nil, fmt.Errorf("comment cache is nil")
	}
	if conn == nil {
		return nil, fmt.Errorf("comment model conn is nil")
	}

	pageKey := ReplyIndexPageKey(
		req.TargetType,
		req.TargetId,
		req.RootId,
		string(req.Sort),
		req.Offset,
		req.Limit,
	)

	if val, err := c.rdb.Get(ctx, pageKey).Result(); err == nil {
		var ids []int64
		if err := json.Unmarshal([]byte(val), &ids); err == nil {
			return ids, nil
		}
	} else if err != redis.Nil {
		return nil, fmt.Errorf("redis get page cache failed, key=%s: %w", pageKey, err)
	}

	sfKey := "reply_ids_page:" + pageKey

	v, err, _ := c.sf.Do(sfKey, func() (interface{}, error) {
		if val, err := c.rdb.Get(ctx, pageKey).Result(); err == nil {
			var ids []int64
			if err := json.Unmarshal([]byte(val), &ids); err == nil {
				return ids, nil
			}
		} else if err != redis.Nil {
			// 双检 redis 异常不拦截，继续回源DB
		}

		ids, err := conn.GetReplyIDsByPage(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("db GetReplyIDsByPage failed: %w", err)
		}

		if b, err := json.Marshal(ids); err == nil {
			_ = c.rdb.Set(ctx, pageKey, b, defaultReplyIDsPageTTL).Err()
		}

		//预加载下一页
		//go c.preloadNextReplyIDsPage(context.Background(), req, conn)

		return ids, nil
	})
	if err != nil {
		return nil, err
	}

	ids, ok := v.([]int64)
	if !ok {
		return nil, fmt.Errorf("singleflight result type assert failed for key=%s", sfKey)
	}
	return ids, nil
}

func (c *CommentCache) preloadNextReplyIDsPage(
	ctx context.Context,
	req model2.GetReplyIDsPageReq,
	conn *model2.CommentModel,
) {
	// 基础保护
	if c == nil || c.rdb == nil || conn == nil {
		return
	}
	if req.Limit <= 0 {
		return
	}

	nextReq := req
	nextReq.Offset = req.Offset + req.Limit

	nextKey := ReplyIndexPageKey(
		nextReq.TargetType,
		nextReq.TargetId,
		nextReq.RootId,
		string(nextReq.Sort),
		nextReq.Offset,
		nextReq.Limit,
	)

	// 已存在就不预加载
	exists, err := c.rdb.Exists(ctx, nextKey).Result()
	if err == nil && exists > 0 {
		return
	}

	// 为避免多个 goroutine 同时预加载，给 nextKey 也走一次 singleflight
	sfKey := "reply_ids_page_preload:" + nextKey
	_, _, _ = c.sf.Do(sfKey, func() (interface{}, error) {
		// 双检
		if val, err := c.rdb.Get(ctx, nextKey).Result(); err == nil && val != "" {
			return nil, nil
		}

		ids, err := conn.GetReplyIDsByPage(ctx, nextReq)
		if err != nil {
			return nil, err
		}
		if b, err := json.Marshal(ids); err == nil {
			_ = c.rdb.Set(ctx, nextKey, b, defaultReplyIDsPageTTL).Err()
		}
		return nil, nil
	})
}
