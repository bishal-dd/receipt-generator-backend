package redisUtil

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)



func DeleteCachePages(cache *redis.Client, ctx context.Context, groupKey string, userId string) error {
	ReceiptsPageGroupKey := fmt.Sprintf("%s:%s", groupKey, userId)
	pageKeys, err := cache.SMembers(ctx, ReceiptsPageGroupKey).Result()
	if err != nil && err != redis.Nil {
        return err
    }
    // Delete all cached pages
    if len(pageKeys) > 0 {
        if err := cache.Del(ctx, pageKeys...).Err(); err != nil {
            return err
        }
        if err := cache.Del(ctx, ReceiptsPageGroupKey).Err(); err != nil {
        return err
    }

    }
    // Clear the set of page keys
    if err := cache.Del(ctx, ReceiptsPageGroupKey).Err(); err != nil {
        return err
    }

	return nil
}