package cache

import (
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

// 方便留接口实现其他缓存
type CommentCache struct {
	rdb *redis.Client
	sf  singleflight.Group
}

func NewCommentCache(rdb *redis.Client) *CommentCache {
	return &CommentCache{
		rdb: rdb,
	}
}
