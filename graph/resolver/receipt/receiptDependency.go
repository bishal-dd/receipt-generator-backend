package receipt

import (
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

type ReceiptResolver struct {
	db            *gorm.DB
	httpClient    *resty.Client
	publicKeyPEM  string
	privateKeyPEM string
}

func InitializeReceiptResolver(db *gorm.DB, httpClient *resty.Client, publicKeyPEM string, privateKeyPEM string) *ReceiptResolver {
	return &ReceiptResolver{
		db:            db,
		httpClient:    httpClient,
		publicKeyPEM:  publicKeyPEM,
		privateKeyPEM: privateKeyPEM,
	}
}

const ReceiptsKey = "receipts"
const ReceiptKey = "receipt"
const ReceiptsPageGroupKey = "receipts:pages"
