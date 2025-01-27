package product

import (
	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ProductResolver struct {
	redis *redis.Client
	db *gorm.DB
	httpClient *resty.Client

}

func InitializeProductResolver(redis *redis.Client, db *gorm.DB, httpClient *resty.Client) *ProductResolver {
	return &ProductResolver{
		db: db,
		redis: redis,
		httpClient: httpClient,
	}
}

const ProductsKey = "products"
const ProductKey = "product"
