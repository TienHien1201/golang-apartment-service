package service

import (
	"context"
	"time"
)

// CacheService defines the caching contract for the domain and usecase layers.
// Implementations (e.g. RedisCache) live in pkg/cache and are injected via DI.
type CacheService interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
	DeleteByPattern(ctx context.Context, pattern string) error
	Exists(ctx context.Context, key string) (bool, error)
	Increment(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)
	GenerateKeyWithParams(key string, params ...interface{}) string
	ErrCacheMiss() error
}
