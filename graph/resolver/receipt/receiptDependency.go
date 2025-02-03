package receipt

import (
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

type ReceiptResolver struct {
	db *gorm.DB
	httpClient *resty.Client

}

func InitializeReceiptResolver( db *gorm.DB,httpClient *resty.Client ) *ReceiptResolver {
	return &ReceiptResolver{
		db: db,
		httpClient: httpClient,
	}
}

const ReceiptsKey = "receipts"
const ReceiptKey = "receipt"
const ReceiptsPageGroupKey = "receipts:pages"