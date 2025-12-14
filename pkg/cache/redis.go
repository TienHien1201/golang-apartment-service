package xcache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host         string
	Port         int
	Password     string
	DB           int
	PoolSize     int
	PoolTimeout  int
	MinIdleConns int
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(cfg *RedisConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		PoolTimeout:  time.Duration(cfg.PoolTimeout) * time.Second,
		MinIdleConns: cfg.MinIdleConns,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RedisCache{
		client: client,
	}, nil
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, expiration).Err()
}

func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrCacheMiss
		}
		return err
	}
	return json.Unmarshal(data, dest)
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *RedisCache) DeleteByPattern(ctx context.Context, pattern string) error {
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	return c.client.Del(ctx, keys...).Err()
}

func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func (c *RedisCache) Increment(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

func (c *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return c.client.Expire(ctx, key, expiration).Result()
}

func (c *RedisCache) GenerateKeyWithParams(key string, params ...interface{}) string {

	return fmt.Sprintf(key, params...)
}

func (c *RedisCache) ErrCacheMiss() error {
	return ErrCacheMiss
}
