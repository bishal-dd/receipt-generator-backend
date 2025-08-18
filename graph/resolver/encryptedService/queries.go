package encryptedservice

import (
	"context"
	"fmt"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
)

func (r *EncryptedServiceResolver) EncryptedServiceByReceiptID(ctx context.Context, receiptId string) ([]*model.EncryptedService, error) {
	var services []*model.EncryptedService
	if err := r.db.Where("encrypted_receipt_id = ?", receiptId).Find(&services).Error; err != nil {
		return nil, err
	}

	// 🔐 Decrypt each receipt
	for _, service := range services {
		err := r.decryptService(service)
		if err != nil {
			// Optionally skip or return error
			return nil, fmt.Errorf("failed to decrypt receipt %s: %w", service.ID, err)
		}
	}

	return services, nil
}

func (r *EncryptedServiceResolver) EncryptedService(ctx context.Context, id string) (*model.EncryptedService, error) {
	service, err := r.GetEncryptedServiceFromDB(id)
	if err != nil {
		return nil, err
	}
	return service, nil
}
