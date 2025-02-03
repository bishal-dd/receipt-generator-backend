package service

import (
	"gorm.io/gorm"
)

type ServiceResolver struct {
	db *gorm.DB
}

func InitializeServiceResolver( db *gorm.DB) *ServiceResolver {
	return &ServiceResolver{
		db: db,
	}
}

const ServicesKey = "services"
const ServiceKey = "service"
