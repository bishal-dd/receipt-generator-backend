package resolver

import (
	"github.com/bishal-dd/receipt-generator-backend/graph/resolver/profile"
	"github.com/bishal-dd/receipt-generator-backend/graph/resolver/receipt"
	"github.com/bishal-dd/receipt-generator-backend/graph/resolver/receiptPDFGenerator"
	"github.com/bishal-dd/receipt-generator-backend/graph/resolver/service"
	"github.com/bishal-dd/receipt-generator-backend/graph/resolver/user"
	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	*user.UserResolver
	*receipt.ReceiptResolver
	*profile.ProfileResolver
	*service.ServiceResolver
	*receiptPDFGenerator.ReceiptPDFGeneratorResolver
}

func InitializeResolver(redis *redis.Client, db *gorm.DB, httpClient *resty.Client) *Resolver {
	return &Resolver{
		UserResolver:    user.InitializeUserResolver(redis, db),
		ReceiptResolver: receipt.InitializeReceiptResolver(redis, db, httpClient),
		ProfileResolver: profile.InitializeProfileResolver(redis, db),
		ServiceResolver: service.InitializeServiceResolver(redis, db),
		ReceiptPDFGeneratorResolver: receiptPDFGenerator.InitializeReceiptPDFGeneratorResolver(redis, db, httpClient ),
	}
}