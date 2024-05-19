package redis

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func Init() *redis.Client {
    redisURL := os.Getenv("REDIS_URL")
    options, err := redis.ParseURL(redisURL)
    if err != nil {
        log.Fatalf("Could not parse Redis URL: %v", err)
    }

    rdb := redis.NewClient(options)

    // Check if the connection is established
    _, err = rdb.Ping(ctx).Result()
    if err != nil {
        log.Fatalf("Could not connect to Redis: %v", err)
    }else{
		log.Println("Successfully connected to Redis")
	}
	
    return rdb
}
