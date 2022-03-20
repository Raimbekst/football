package redis

import (
	"carWash/internal/config"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.TODO()
)

func NewRedisDB(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("redis.NewRedisDB: %w", err)
	}
	return client, nil

}
