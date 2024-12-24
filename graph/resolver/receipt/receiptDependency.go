package receipt

import (
	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ReceiptResolver struct {
	redis *redis.Client
	db *gorm.DB
	httpClient *resty.Client

}

func InitializeReceiptResolver(redis *redis.Client, db *gorm.DB,httpClient *resty.Client ) *ReceiptResolver {
	return &ReceiptResolver{
		db: db,
		redis: redis,
		httpClient: httpClient,
	}
}

const ReceiptsKey = "receipts"
const ReceiptKey = "receipt"
const ReceiptsPageGroupKey = "receipts:pages"