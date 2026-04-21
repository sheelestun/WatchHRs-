package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	rdb *redis.Client
}

func NewRedisCache(rdb *redis.Client) *RedisCache {
	return &RedisCache{rdb: rdb}
}

func (r *RedisCache) SaveTokenInCache(ctx context.Context, tokenID, userID string, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return errors.New("invalid expiration time")
	}

	key := "refresh:" + tokenID

	return r.rdb.Set(ctx, key, userID, ttl).Err()
}

func (r *RedisCache) ExistsTokenInCache(ctx context.Context, tokenID string) (bool, error) {
	key := "refresh:" + tokenID

	n, err := r.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return n == 1, nil
}

func (r *RedisCache) DeleteTokenInCache(ctx context.Context, tokenID string) error {
	key := "refresh:" + tokenID

	return r.rdb.Del(ctx, key).Err()
}
