package redisdb

import (
	"context"
	"github.com/Turalchik/authentication-service/internal/apperrors"
	"github.com/go-redis/redis/v8"
	"time"
)

func NewRedisClient(addr, password string, db int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Пингуем Redis, чтобы убедиться, что он доступен
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, apperrors.ErrRedisPingFailed
	}

	return rdb, nil
}
