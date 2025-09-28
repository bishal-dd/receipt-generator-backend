package receiptFiles

import (
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

type ReceiptFilesResolver struct {
	db            *gorm.DB
	httpClient    *resty.Client
	publicKeyPEM  string
	privateKeyPEM string
}

func InitializeReceiptResolver(db *gorm.DB, httpClient *resty.Client, publicKeyPEM string, privateKeyPEM string) *ReceiptFilesResolver {
	return &ReceiptFilesResolver{
		db:            db,
		httpClient:    httpClient,
		publicKeyPEM:  publicKeyPEM,
		privateKeyPEM: privateKeyPEM,
	}
}

const ReceiptsKey = "receipts"
const ReceiptKey = "receipt"
const ReceiptsPageGroupKey = "receipts:pages"
