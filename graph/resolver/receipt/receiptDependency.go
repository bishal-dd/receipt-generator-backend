package receipt

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ReceiptResolver struct {
	redis *redis.Client
	db *gorm.DB
}

func InitializeReceiptResolver(redis *redis.Client, db *gorm.DB) *ReceiptResolver {
	return &ReceiptResolver{
		db: db,
		redis: redis,
	}
}

const ReceiptsKey = "receipts"
const ReceiptKey = "receipt"
const ReceiptsPageGroupKey = "receipts:pages"