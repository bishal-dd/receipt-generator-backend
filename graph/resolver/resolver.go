package resolver

import (
	"github.com/bishal-dd/receipt-generator-backend/graph/resolver/receipt"
	"github.com/bishal-dd/receipt-generator-backend/graph/resolver/user"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	*user.UserResolver
	*receipt.ReceiptResolver
}

func InitializeResolver(redis *redis.Client, db *gorm.DB) *Resolver {
	return &Resolver{
		UserResolver:    user.InitializeUserResolver(redis, db),
		ReceiptResolver: receipt.InitializeReceiptResolver(redis, db),
	}
}