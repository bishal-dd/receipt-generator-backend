package redisUtil

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func DeleteCacheItem(redis *redis.Client, ctx context.Context, key string, id string ) error {
	return redis.Del(ctx, fmt.Sprintf("%s:%s", key, id)).Err();
}