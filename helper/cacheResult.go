package helper

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func CacheResult(redis *redis.Client, ctx context.Context, key string, value interface{}, ttl int ) error {
	valueJSON, err := Marshal(value)
	if err != nil {
		return  err
	}
	return redis.Set(ctx, key, valueJSON, time.Duration(ttl)*time.Minute).Err();
}