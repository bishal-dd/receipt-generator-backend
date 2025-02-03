package receiptPDFGenerator

import (
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

type ReceiptPDFGeneratorResolver struct {
	db *gorm.DB
	httpClient *resty.Client
}

func InitializeReceiptPDFGeneratorResolver( db *gorm.DB, httpClient *resty.Client) *ReceiptPDFGeneratorResolver {
	return &ReceiptPDFGeneratorResolver{
		db: db,
		httpClient: httpClient,
	}
}
