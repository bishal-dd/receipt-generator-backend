package encryptedServiceLoader

import (
	"context"
	"fmt"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"gorm.io/gorm"
)

type EncryptedServiceReader struct {
	db *gorm.DB
}

func NewEncryptedServiceReader(db *gorm.DB) *EncryptedServiceReader {
	return &EncryptedServiceReader{db: db}
}
func (u *EncryptedServiceReader) GetEncryptedServicesByReceiptIds(ctx context.Context, receiptIds []string) ([][]*model.EncryptedService, []error) {
	if len(receiptIds) == 0 {
		return [][]*model.EncryptedService{}, nil
	}

	var encryptedServices []*model.EncryptedService
	err := u.db.Where("encrypted_receipt_id IN (?)", receiptIds).Find(&encryptedServices).Error
	if err != nil {
		return nil, []error{fmt.Errorf("failed to fetch encryptedServices: %w", err)}
	}

	encryptedServiceMap := make(map[string][]*model.EncryptedService)
	for _, encryptedService := range encryptedServices {
		encryptedServiceMap[encryptedService.EncryptedReceiptID] = append(encryptedServiceMap[encryptedService.EncryptedReceiptID], encryptedService)
	}

	result := make([][]*model.EncryptedService, len(receiptIds))
	for i, userID := range receiptIds {
		result[i] = encryptedServiceMap[userID]
	}

	return result, nil
}
