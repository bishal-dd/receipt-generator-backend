package resolver

import (
	"github.com/bishal-dd/receipt-generator-backend/graph/resolver/encryptedReceipt"
	encryptedservice "github.com/bishal-dd/receipt-generator-backend/graph/resolver/encryptedService"
	"github.com/bishal-dd/receipt-generator-backend/graph/resolver/product"
	"github.com/bishal-dd/receipt-generator-backend/graph/resolver/profile"
	"github.com/bishal-dd/receipt-generator-backend/graph/resolver/receipt"
	"github.com/bishal-dd/receipt-generator-backend/graph/resolver/receiptPDFGenerator"
	"github.com/bishal-dd/receipt-generator-backend/graph/resolver/service"
	"github.com/bishal-dd/receipt-generator-backend/graph/resolver/user"
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	*user.UserResolver
	*receipt.ReceiptResolver
	*encryptedReceipt.EncryptedReceiptResolver
	*encryptedservice.EncryptedServiceResolver
	*profile.ProfileResolver
	*service.ServiceResolver
	*receiptPDFGenerator.ReceiptPDFGeneratorResolver
	*product.ProductResolver
}

func InitializeResolver(db *gorm.DB, httpClient *resty.Client, publicKeyPEM string, privateKeyPEM string) *Resolver {
	return &Resolver{
		UserResolver:                user.InitializeUserResolver(db),
		ReceiptResolver:             receipt.InitializeReceiptResolver(db, httpClient),
		ProfileResolver:             profile.InitializeProfileResolver(db),
		ServiceResolver:             service.InitializeServiceResolver(db),
		ReceiptPDFGeneratorResolver: receiptPDFGenerator.InitializeReceiptPDFGeneratorResolver(db, httpClient, publicKeyPEM, privateKeyPEM),
		ProductResolver:             product.InitializeProductResolver(db, httpClient),
		EncryptedReceiptResolver:    encryptedReceipt.InitializeEncryptedReceiptResolver(db, httpClient, publicKeyPEM, privateKeyPEM),
		EncryptedServiceResolver:    encryptedservice.InitializeEncryptedServiceResolver(db, httpClient, publicKeyPEM, privateKeyPEM),
	}
}
