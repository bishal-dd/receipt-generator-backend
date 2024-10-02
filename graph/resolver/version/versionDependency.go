package version

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type VersionResolver struct {
	redis *redis.Client
	db *gorm.DB
}

func InitializeVersionResolver(redis *redis.Client, db *gorm.DB) *VersionResolver {
	return &VersionResolver{
		db: db,
		redis: redis,
	}
}

const VersionsKey = "versions"
const VersionKey = "version"
