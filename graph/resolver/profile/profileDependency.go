package profile

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ProfileResolver struct {
	redis *redis.Client
	db *gorm.DB
}

func InitializeProfileResolver(redis *redis.Client, db *gorm.DB) *ProfileResolver {
	return &ProfileResolver{
		db: db,
		redis: redis,
	}
}

const ProfilesKey = "profiles"
const ProfileKey = "profile:"
