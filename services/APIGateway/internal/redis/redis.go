package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sheelestun/WatchHRs-/internal/config"
	log "github.com/sirupsen/logrus"
)

func Connect(cfg *config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Info("Successfully connected to Redis")

	return rdb, nil
}
