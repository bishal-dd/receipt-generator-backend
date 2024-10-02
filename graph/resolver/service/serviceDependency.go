package service

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ServiceResolver struct {
	redis *redis.Client
	db *gorm.DB
}

func InitializeServiceResolver(redis *redis.Client, db *gorm.DB) *ServiceResolver {
	return &ServiceResolver{
		db: db,
		redis: redis,
	}
}

const ServicesKey = "services"
const ServiceKey = "service"
