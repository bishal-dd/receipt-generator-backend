package redisUtil

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func CacheResultString(redis *redis.Client, ctx context.Context, key string, value string, ttl int ) error {
	return redis.Set(ctx, key, value, time.Duration(ttl)*time.Minute).Err();
}