package utils

import (
	"context"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
)

func TestRedisCache_SetGetDeletePattern(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis start failed: %v", err)
	}
	defer s.Close()

	cfg := CacheConfig{
		Addr:     s.Addr(),
		Password: "",
		TTL:      2,
		Enabled:  true,
	}

	rc, err := NewRedisCache(cfg)
	if err != nil {
		t.Fatalf("NewRedisCache failed: %v", err)
	}
	defer rc.Close()

	ctx := context.Background()

	key := "favorite:user:1:folder:2:page:1"
	payload := map[string]interface{}{"items": []int{1, 2, 3}, "total": 3}

	if err := rc.Set(ctx, key, payload); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	var got map[string]interface{}
	if err := rc.Get(ctx, key, &got); err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got == nil {
		t.Fatalf("expected non-nil result from Get")
	}

	// pattern delete
	pattern := "favorite:user:1:folder:2:*"
	if err := rc.DeletePattern(ctx, pattern); err != nil {
		t.Fatalf("DeletePattern failed: %v", err)
	}

	var after map[string]interface{}
	if err := rc.Get(ctx, key, &after); err == nil && after != nil {
		t.Fatalf("expected key to be deleted by DeletePattern")
	}

	// test TTL expiration
	if err := rc.Set(ctx, key, payload, time.Duration(1)*time.Second); err != nil {
		t.Fatalf("Set with ttl failed: %v", err)
	}
	// fast-forward miniredis time to ensure expiration
	s.FastForward(2 * time.Second)
	var expired map[string]interface{}
	if err := rc.Get(ctx, key, &expired); err == nil && expired != nil {
		t.Fatalf("expected key to expire")
	}
}
