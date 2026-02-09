package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheConfig Redis 缓存配置
type CacheConfig struct {
	// Redis 连接地址
	Addr string
	// Redis 密码
	Password string
	// 缓存过期时间（秒）
	TTL int
	// 是否启用缓存
	Enabled bool
}

// RedisCache Redis 缓存实现
type RedisCache struct {
	client  *redis.Client
	ttl     time.Duration
	enabled bool
}

// CacheKey 缓存键生成器
type CacheKey struct {
	prefix string
}

// NewRedisCache 创建 Redis 缓存客户端
func NewRedisCache(config CacheConfig) (*RedisCache, error) {
	if !config.Enabled {
		return &RedisCache{
			enabled: false,
		}, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           0,
		MaxRetries:   3,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return &RedisCache{
		client:  client,
		ttl:     time.Duration(config.TTL) * time.Second,
		enabled: true,
	}, nil
}

// Get 从缓存中获取数据
func (rc *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	if !rc.enabled || rc.client == nil {
		return fmt.Errorf("cache disabled or not initialized")
	}

	val, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil // 缓存未找到，返回 nil
		}
		return fmt.Errorf("redis get failed: %w", err)
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("json unmarshal failed: %w", err)
	}

	return nil
}

// Set 将数据存入缓存
func (rc *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl ...interface{}) error {
	if !rc.enabled || rc.client == nil {
		return nil // 缓存禁用时不报错，只忽略
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}

	expiration := rc.ttl
	if len(ttl) > 0 {
		// 尝试处理不同的 ttl 类型
		if t, ok := ttl[0].(time.Duration); ok {
			expiration = t
		}
	}

	if err := rc.client.Set(ctx, key, string(data), expiration).Err(); err != nil {
		return fmt.Errorf("redis set failed: %w", err)
	}

	return nil
}

// Delete 删除缓存
func (rc *RedisCache) Delete(ctx context.Context, keys ...string) error {
	if !rc.enabled || rc.client == nil {
		return nil
	}

	if len(keys) == 0 {
		return nil
	}

	if err := rc.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("redis delete failed: %w", err)
	}

	return nil
}

// DeletePattern 删除匹配模式的所有缓存
// 例如: favorite:user:1:* 会删除所有该用户的缓存
func (rc *RedisCache) DeletePattern(ctx context.Context, pattern string) error {
	if !rc.enabled || rc.client == nil {
		return nil
	}

	// 使用 KEYS 命令查找匹配的 key（注意：在生产环境中应使用 SCAN）
	keys, err := rc.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("redis keys failed: %w", err)
	}

	if len(keys) > 0 {
		if err := rc.client.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("redis delete failed: %w", err)
		}
	}

	return nil
}

// Exists 检查缓存是否存在
func (rc *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	if !rc.enabled || rc.client == nil {
		return false, nil
	}

	exists, err := rc.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis exists failed: %w", err)
	}

	return exists > 0, nil
}

// Incr 增加计数器
func (rc *RedisCache) Incr(ctx context.Context, key string) (int64, error) {
	if !rc.enabled || rc.client == nil {
		return 0, fmt.Errorf("cache disabled")
	}

	val, err := rc.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis incr failed: %w", err)
	}

	return val, nil
}

// Close 关闭 Redis 连接
func (rc *RedisCache) Close() error {
	if rc.client == nil {
		return nil
	}
	return rc.client.Close()
}
