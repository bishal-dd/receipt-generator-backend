package encryptedReceipt

import (
	"context"
	"fmt"

	"github.com/bishal-dd/receipt-generator-backend/graph/loaders"
	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/encryption"
	"github.com/bishal-dd/receipt-generator-backend/helper/stringUtil"
)

// func (r *EncryptedReceiptResolver) LoadServiceFromEncryptedReceipts(ctx context.Context, receipts []*model.EncryptedReceipt) ([]*model.EncryptedReceipt, error) {
// 	loaders := loaders.For(ctx)
// 	receiptIds := make([]string, len(receipts))
// 	for i, receipt := range receipts {
// 		receiptIds[i] = receipt.ID
// 	}

// 	serviceResults, err := loaders.ServiceLoader.LoadAll(ctx, receiptIds)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for i, services := range serviceResults {
// 		receipts[i].Services = services
// 	}
// 	return receipts, nil
// }

// func (r *EncryptedReceiptResolver) LoadServiceFromEncryptedReceipt(ctx context.Context, receipt *model.EncryptedReceipt) (*model.EncryptedReceipt, error) {
// 	loaders := loaders.For(ctx)
// 	serviceResults, err := loaders.ServiceLoader.Load(ctx, receipt.ID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	receipt.Services = serviceResults

// 	return receipt, nil
// }

func (r *EncryptedReceiptResolver) LoadServiceFromEncryptedReceipt(ctx context.Context, receipt *model.EncryptedReceipt) (*model.EncryptedReceipt, error) {
	loaders := loaders.For(ctx)
	serviceResults, err := loaders.EncryptedServiceLoader.Load(ctx, receipt.ID)
	if err != nil {
		return nil, err
	}
	// 🔐 Decrypt each receipt
	for _, service := range serviceResults {
		err := r.decryptService(service)
		if err != nil {
			// Optionally skip or return error
			return nil, fmt.Errorf("failed to decrypt receipt %s: %w", service.ID, err)
		}
	}

	receipt.EncryptedServices = serviceResults

	return receipt, nil
}

func (r *EncryptedReceiptResolver) decryptService(service *model.EncryptedService) error {
	aesKey, iv, err := encryption.DecryptKeyAndIV(r.privateKeyPEM, *service.AesKeyEncrypted, *service.AesIv)
	if err != nil {
		return fmt.Errorf("decrypt AES key: %w", err)
	}
	// 4. Apply decryption
	service.Description = stringUtil.DerefString(encryption.DecryptField(stringUtil.StrPtr(service.Description), aesKey, iv))
	service.Amount = stringUtil.DerefString(encryption.DecryptField(stringUtil.StrPtr(service.Amount), aesKey, iv))
	service.Quantity = stringUtil.DerefString(encryption.DecryptField(stringUtil.StrPtr(service.Quantity), aesKey, iv))
	service.Rate = stringUtil.DerefString(encryption.DecryptField(stringUtil.StrPtr(service.Rate), aesKey, iv))

	return nil
}
