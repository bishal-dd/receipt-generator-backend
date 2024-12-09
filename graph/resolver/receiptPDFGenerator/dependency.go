package receiptPDFGenerator

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ReceiptPDFGeneratorResolver struct {
	redis *redis.Client
	db *gorm.DB
}

func InitializeReceiptPDFGeneratorResolver(redis *redis.Client, db *gorm.DB) *ReceiptPDFGeneratorResolver {
	return &ReceiptPDFGeneratorResolver{
		db: db,
		redis: redis,
	}
}
