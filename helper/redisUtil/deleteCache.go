package redisutil

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/helper"
	"github.com/redis/go-redis/v9"
)


func DeleteCache[T any](cachedValues []T, key string, id string, jsonValue string, ctx context.Context, redis *redis.Client )  error  {
	if err := helper.Unmarshal([]byte(jsonValue), &cachedValues); err != nil {
		return err
	}
	// Find the index of the deleted receipt in the cached list
	index := -1
	for i, value := range cachedValues {
		if value.ID == id {
			index = i
			break
		}
	}
	// If the receipt is found, remove it from the cached list
	if index != -1 {
		cachedValues = append(cachedValues[:index], cachedValues[index+1:]...)
		// Update the cache with the modified list of receipts
		if err := helper.CacheResult(redis, ctx, key, cachedValues, 10); err != nil {
			return err
		}
	}

	return nil
}