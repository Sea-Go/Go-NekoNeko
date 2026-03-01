package cache

import (
	"context"
	"encoding/json"
	"fmt"
	model2 "sea-try-go/service/comment/rpc/internal/model"
	"time"
)

const defaultReplyContentTTL = 24 * time.Hour

func (cache *CommentCache) BatchGetContentCache(ctx context.Context, ids []int64, conn *model2.CommentModel) ([]model2.CommentContent, error) {
	if cache == nil || cache.rdb == nil {
		return nil, fmt.Errorf("comment cache is nil")
	}
	if conn == nil {
		return nil, fmt.Errorf("comment model conn is nil")
	}
	if len(ids) == 0 {
		return []model2.CommentContent{}, nil
	}

	//过滤ID
	validIDs := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id > 0 {
			validIDs = append(validIDs, id)
		}
	}
	if len(validIDs) == 0 {
		return []model2.CommentContent{}, nil
	}

	keys := make([]string, 0, len(validIDs))
	for _, id := range validIDs {
		keys = append(keys, ReplyContentKey(id))
	}

	values, err := cache.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("redis MGet content cache failed: %w", err)
	}

	contentMap := make(map[int64]model2.CommentContent, len(validIDs))
	missIDs := make([]int64, 0)

	// 4) 解析 MGET 结果（顺序与 keys / validIDs 对应）
	for i, v := range values {
		id := validIDs[i]

		if v == nil {
			missIDs = append(missIDs, id)
			continue
		}

		// go-redis MGet 返回 []interface{}，常见是 string 或 []byte
		var raw []byte
		switch vv := v.(type) {
		case string:
			raw = []byte(vv)
		case []byte:
			raw = vv
		default:
			// 类型异常，按 miss 处理更稳
			missIDs = append(missIDs, id)
			continue
		}

		var c model2.CommentContent
		if err := json.Unmarshal(raw, &c); err != nil {
			missIDs = append(missIDs, id)
			continue
		}
		contentMap[id] = c
	}

	// 5) miss 回源 DB
	if len(missIDs) > 0 {
		dbContents, err := conn.BatchGetReplyContent(ctx, missIDs)
		if err != nil {
			return nil, fmt.Errorf("db fallback BatchGetReplyContentByCommentIDs failed: %w", err)
		}

		// 回填 map（供最终按顺序输出）
		for _, c := range dbContents {
			contentMap[c.CommentId] = c
		}

		// 6) 回填 Redis（pipeline）
		if len(dbContents) > 0 {
			pipe := cache.rdb.Pipeline()
			for _, c := range dbContents {
				b, err := json.Marshal(c)
				if err != nil {
					continue
				}
				pipe.Set(ctx, ReplyContentKey(c.CommentId), b, defaultReplyContentTTL)
			}
			_, _ = pipe.Exec(ctx)
		}
	}

	// 7) 按输入 ids 顺序重排输出（关键）
	result := make([]model2.CommentContent, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if c, ok := contentMap[id]; ok {
			result = append(result, c)
		}
	}

	return result, nil
}
