package redisUtil

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)


func CachePages(redis *redis.Client, groupKey string, ctx context.Context, key string, value interface{}, offset int, limit int, userId string) error {
	pageCacheKey := fmt.Sprintf("%s:%d:%d:%s", key, offset, limit,userId)
    pagesGroupKey := fmt.Sprintf("%s:%s",groupKey, userId )
	if err := redis.SAdd(ctx, pagesGroupKey, pageCacheKey).Err(); err != nil {
		return err
	}
	if err := redis.Expire(ctx, pagesGroupKey, 600*time.Second).Err(); err != nil {
		return err
	}
	if err := CacheResult(redis, ctx, pageCacheKey, value, 10); err != nil {
		return err
	}

	return nil
}