package user

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserResolver struct{
	redis *redis.Client
	db *gorm.DB
}

func InitializeUserResolver(redis *redis.Client, db *gorm.DB) *UserResolver {
	return &UserResolver{
		db: db,
		redis: redis,
	}
}