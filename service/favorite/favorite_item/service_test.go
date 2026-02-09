package favorite_item

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// fakeRepo 是一个内存实现，用于测试 Service 缓存行为
type fakeRepo struct {
	items []*FavoriteItem
	calls int
}

func (f *fakeRepo) Add(ctx context.Context, item *FavoriteItem) error {
	item.ID = int64(len(f.items) + 1)
	f.items = append([]*FavoriteItem{item}, f.items...)
	return nil
}
func (f *fakeRepo) Delete(ctx context.Context, userID int64, objectType string, objectID int64) error {
	return nil
}
func (f *fakeRepo) ListByFolder(ctx context.Context, folderID int64, offset, limit int) ([]*FavoriteItem, error) {
	f.calls++
	// return a copy
	var res []*FavoriteItem
	for _, it := range f.items {
		if it.FolderID == folderID {
			res = append(res, it)
		}
	}
	if offset >= len(res) {
		return []*FavoriteItem{}, nil
	}
	end := offset + limit
	if end > len(res) {
		end = len(res)
	}
	return res[offset:end], nil
}
func (f *fakeRepo) CheckExists(ctx context.Context, userID int64, objectType string, objectID int64) (bool, int64, int64, error) {
	return false, 0, 0, nil
}
func (f *fakeRepo) CountByFolder(ctx context.Context, folderID int64) (int64, error) {
	var c int64
	for _, it := range f.items {
		if it.FolderID == folderID {
			c++
		}
	}
	return c, nil
}

func TestService_ListItems_CacheHitInvalidate(t *testing.T) {
	// miniredis
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis start failed: %v", err)
	}
	defer s.Close()

	// create fake repo and seed one item
	repo := &fakeRepo{}
	repo.items = []*FavoriteItem{
		{ID: 1, UserID: 1, FolderID: 2, ObjectType: "url", ObjectID: 100},
	}

	// create service directly
	svc := &Service{
		itemRepo: repo,
		db:       &pgxpool.Pool{},
		cache:    &NopCache{},
		cacheKey: NewCacheKeyBuilder(),
	}

	// replace cache with real redis client backed implementation for tests
	rcli := redis.NewClient(&redis.Options{Addr: s.Addr(), Password: "", DB: 0})
	testCache := &redisCacheTest{client: rcli}
	defer rcli.Close()
	svc.cache = testCache

	ctx := context.Background()

	// first call: should hit repo (calls==1)
	_, _, err = svc.ListItems(ctx, 1, 2, 1, 10)
	if err != nil {
		t.Fatalf("ListItems failed: %v", err)
	}
	if repo.calls != 1 {
		t.Fatalf("expected repo.calls==1, got %d", repo.calls)
	}

	// second call: should be served from cache (calls still 1)
	_, _, err = svc.ListItems(ctx, 1, 2, 1, 10)
	if err != nil {
		t.Fatalf("ListItems failed: %v", err)
	}
	if repo.calls != 1 {
		t.Fatalf("expected repo.calls still 1 after cache hit, got %d", repo.calls)
	}

	// create a new item to invalidate cache
	if _, err := svc.CreateItem(ctx, 1, 2, "url", 200); err != nil {
		t.Fatalf("CreateItem failed: %v", err)
	}

	// after invalidation, listing should hit repo again
	_, _, err = svc.ListItems(ctx, 1, 2, 1, 10)
	if err != nil {
		t.Fatalf("ListItems failed: %v", err)
	}
	if repo.calls != 2 {
		t.Fatalf("expected repo.calls==2 after invalidation, got %d", repo.calls)
	}

	// cleanup deferred
}

// redisCacheTest 是测试用的 CacheInterface 实现，使用 go-redis
type redisCacheTest struct {
	client *redis.Client
}

func (r *redisCacheTest) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (r *redisCacheTest) Set(ctx context.Context, key string, value interface{}, ttl ...interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, string(data), 0).Err()
}

func (r *redisCacheTest) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *redisCacheTest) DeletePattern(ctx context.Context, pattern string) error {
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}
	return nil
}

func (r *redisCacheTest) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}
