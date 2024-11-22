package redis

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

// Redis clients for different purposes
var (
	CacheClient *redis.Client
	QueueClient *redis.Client
)

// Config holds Redis configuration
type Config struct {
	URL      string
	Database int // Different database numbers for separation
}

func createClient(cfg Config) (*redis.Client, error) {
	options, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, err
	}
	
	// Override the database number
	options.DB = cfg.Database

	client := redis.NewClient(options)

	// Check if the connection is established
	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return client, nil
}

// Init initializes both Redis clients
func Init() (*redis.Client, *redis.Client, error) {
	// Configuration for cache client (using database 0)
	cacheConfig := Config{
		URL:      os.Getenv("REDIS_CACHE_URL"),
		Database: 0,
	}

	// Configuration for queue client (using database 1)
	queueConfig := Config{
		URL:      os.Getenv("REDIS_QUEUE_URL"), // You can use a different URL if needed
		Database: 1,
	}

	// Initialize cache client
	var err error
	CacheClient, err = createClient(cacheConfig)
	if err != nil {
		return nil,nil,err
	}
	log.Println("Successfully connected to Redis Cache (DB: 0)")

	// Initialize queue client
	QueueClient, err = createClient(queueConfig)
	if err != nil {
		return nil, nil,err
	}
	log.Println("Successfully connected to Redis Queue (DB: 1)")

	return CacheClient, QueueClient,nil
}

// Close closes both Redis clients
func Close() {
	if CacheClient != nil {
		if err := CacheClient.Close(); err != nil {
			log.Printf("Error closing cache client: %v", err)
		}
	}
	if QueueClient != nil {
		if err := QueueClient.Close(); err != nil {
			log.Printf("Error closing queue client: %v", err)
		}
	}
}