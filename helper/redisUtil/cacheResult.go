package redisUtil

import (
	"context"
	"time"

	"github.com/bishal-dd/receipt-generator-backend/helper/json"
	"github.com/redis/go-redis/v9"
)

func CacheResult(redis *redis.Client, ctx context.Context, key string, value interface{}, ttl int ) error {
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return  err
	}
	return redis.Set(ctx, key, valueJSON, time.Duration(ttl)*time.Minute).Err();
}