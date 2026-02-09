package favorite_item

import (
	"context"
	"fmt"
)

// CacheInterface 缓存接口（避免循环导入）
type CacheInterface interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl ...interface{}) error
	Delete(ctx context.Context, keys ...string) error
	DeletePattern(ctx context.Context, pattern string) error
	Exists(ctx context.Context, key string) (bool, error)
}

// CacheKeyBuilder 缓存键生成工具
type CacheKeyBuilder struct{}

// NewCacheKeyBuilder 创建缓存键生成器
func NewCacheKeyBuilder() *CacheKeyBuilder {
	return &CacheKeyBuilder{}
}

// FavoriteListKey 生成收藏列表缓存键
// 格式: favorite:user:{userID}:folder:{folderID}:page:{page}
func (ckb *CacheKeyBuilder) FavoriteListKey(userID, folderID int64, page int) string {
	return fmt.Sprintf("favorite:user:%d:folder:%d:page:%d", userID, folderID, page)
}

// FavoriteFolderPattern 生成收藏夹所有缓存的匹配模式
// 格式: favorite:user:{userID}:folder:{folderID}:*
func (ckb *CacheKeyBuilder) FavoriteFolderPattern(userID, folderID int64) string {
	return fmt.Sprintf("favorite:user:%d:folder:%d:*", userID, folderID)
}

// NopCache 空缓存实现（当缓存禁用时使用）
type NopCache struct{}

// Get 空实现
func (nc *NopCache) Get(ctx context.Context, key string, dest interface{}) error {
	return nil
}

// Set 空实现
func (nc *NopCache) Set(ctx context.Context, key string, value interface{}, ttl ...interface{}) error {
	return nil
}

// Delete 空实现
func (nc *NopCache) Delete(ctx context.Context, keys ...string) error {
	return nil
}

// DeletePattern 空实现
func (nc *NopCache) DeletePattern(ctx context.Context, pattern string) error {
	return nil
}

// Exists 空实现
func (nc *NopCache) Exists(ctx context.Context, key string) (bool, error) {
	return false, nil
}
