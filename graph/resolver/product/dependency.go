package product

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ProductResolver struct {
	redis *redis.Client
	db *gorm.DB
}

func InitializeProductResolver(redis *redis.Client, db *gorm.DB) *ProductResolver {
	return &ProductResolver{
		db: db,
		redis: redis,
	}
}

const ProductsKey = "products"
const ProductKey = "product"
