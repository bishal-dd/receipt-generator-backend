package receiptPDFGenerator

import (
	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ReceiptPDFGeneratorResolver struct {
	redis *redis.Client
	db *gorm.DB
	httpClient *resty.Client
}

func InitializeReceiptPDFGeneratorResolver(redis *redis.Client, db *gorm.DB, httpClient *resty.Client) *ReceiptPDFGeneratorResolver {
	return &ReceiptPDFGeneratorResolver{
		db: db,
		redis: redis,
		httpClient: httpClient,
	}
}
