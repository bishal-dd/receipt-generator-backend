package encryptedReceipt

import (
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

type EncryptedReceiptResolver struct {
	db            *gorm.DB
	httpClient    *resty.Client
	publicKeyPEM  string
	privateKeyPEM string
}

func InitializeEncryptedReceiptResolver(db *gorm.DB, httpClient *resty.Client, publicKeyPEM string, privateKeyPEM string) *EncryptedReceiptResolver {
	return &EncryptedReceiptResolver{
		db:            db,
		httpClient:    httpClient,
		publicKeyPEM:  publicKeyPEM,
		privateKeyPEM: privateKeyPEM,
	}
}

const EncryptedReceiptsKey = "receipts"
const EncryptedReceiptKey = "receipt"
const EncryptedReceiptsPageGroupKey = "receipts:pages"
