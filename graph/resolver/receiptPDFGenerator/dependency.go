package receiptPDFGenerator

import (
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

type ReceiptPDFGeneratorResolver struct {
	db            *gorm.DB
	httpClient    *resty.Client
	publicKeyPEM  string
	privateKeyPEM string
}

func InitializeReceiptPDFGeneratorResolver(db *gorm.DB, httpClient *resty.Client, publicKeyPEM string, privateKeyPEM string) *ReceiptPDFGeneratorResolver {
	return &ReceiptPDFGeneratorResolver{
		db:            db,
		httpClient:    httpClient,
		publicKeyPEM:  publicKeyPEM,
		privateKeyPEM: privateKeyPEM,
	}
}
