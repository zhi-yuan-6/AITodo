package util

import (
	"AITodo/config"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
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

// 存储refreshtoken
func StoreToken(userID uint, token string, expiration time.Duration) error {
	return Client.Set(ctx, strconv.FormatUint(uint64(userID), 10), token, expiration).Err()
}

func GetToken(userID uint) (string, error) {
	return Client.Get(ctx, strconv.FormatUint(uint64(userID), 10)).Result()
}

// 使令牌失效（加入黑名单）
func InvalidateToken(tokenString string) error {
	// 解析 Token 获取 jti
	claims, err := ParseJWT(tokenString)
	if err != nil {
		return err
	}
	jti := claims.RegisteredClaims.ID // 假设 JWT 包含 jti 字段
	key := fmt.Sprintf("jwt:blacklist:%s", jti)
	return Client.Set(ctx, key, "revoked", time.Until(claims.ExpiresAt.Time)).Err()
}

// 检查令牌是否有效
func IsTokenValid(tokenString string) (bool, error) {
	claims, err := ParseJWT(tokenString)
	if err != nil {
		return false, err
	}

	// 检查黑名单中的 jti
	jti := claims.RegisteredClaims.ID
	exist, err := Client.Exists(ctx, fmt.Sprintf("jwt:blacklist:%s", jti)).Result()
	return exist == 0, err
}
