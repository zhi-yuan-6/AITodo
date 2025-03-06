package util

import (
	"AITodo/config"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	ctx    = context.Background()
	Client *redis.Client
)

func Initialize(cfg config.RedisConfig) error {
	Client = redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     100, // 最大连接数
		MinIdleConns: 20,  // 最小空闲连接
		IdleTimeout:  5 * time.Minute,
	})

	if _, err := Client.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}
	return nil
}

func StoreCode(phone, code string, expiration time.Duration) error {
	return Client.Set(ctx, phone, code, expiration).Err()
}

func GetCode(phone string) (string, error) {
	return Client.Get(ctx, phone).Result()
}
