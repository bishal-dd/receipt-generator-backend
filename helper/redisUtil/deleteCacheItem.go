package redisUtil

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func DeleteCacheItem(r *redis.Client, ctx context.Context, key string, id string ) error {
	return r.Del(ctx, fmt.Sprintf("%s:%s", key, id)).Err();
}