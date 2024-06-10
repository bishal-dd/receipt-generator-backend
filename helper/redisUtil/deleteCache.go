package redisUtil

import (
	"context"
	"errors"
	"reflect"

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
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Struct {
			idField := v.FieldByName("ID")
			if !idField.IsValid() || idField.Kind() != reflect.String {
				return errors.New("ID field not found or not a string")
			}
			if idField.String() == id {
				index = i
				break
			}
		}
	}

	// If the receipt is found, remove it from the cached list
	if index != -1 {
		cachedValues = append(cachedValues[:index], cachedValues[index+1:]...)
		// Update the cache with the modified list of receipts
		if err := CacheResult(redis, ctx, key, cachedValues, 10); err != nil {
			return err
		}
	}

	return nil
}