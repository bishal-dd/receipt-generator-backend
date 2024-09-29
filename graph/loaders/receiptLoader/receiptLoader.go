package receiptLoader

import (
	"context"
	"fmt"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"gorm.io/gorm"
)


type ReceiptReader struct {
	db *gorm.DB
}

func NewReceiptReader(db *gorm.DB) *ReceiptReader {
	return &ReceiptReader{db: db}
}
func (u *ReceiptReader) GetReceiptsByUserIds(ctx context.Context, userIds []string) ([][]*model.Receipt, []error) {
    if len(userIds) == 0 {
        return [][]*model.Receipt{}, nil
    }

    var receipts []*model.Receipt
    err := u.db.Where("user_id IN (?)", userIds).Find(&receipts).Error
    if err != nil {
        return nil, []error{fmt.Errorf("failed to fetch receipts: %w", err)}
    }

    receiptMap := make(map[string][]*model.Receipt)
    for _, receipt := range receipts {
        receiptMap[receipt.UserID] = append(receiptMap[receipt.UserID], receipt)
    }

    result := make([][]*model.Receipt, len(userIds))
    for i, userID := range userIds {
        result[i] = receiptMap[userID] 
    }

    return result, nil
}
