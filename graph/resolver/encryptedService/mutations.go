package encryptedservice

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/encryption"
	"github.com/google/uuid"
)

func (r *EncryptedServiceResolver) CreateEncryptedService(ctx context.Context, input model.CreateEncryptedService) (*model.EncryptedService, error) {
	aesKey, iv, err := encryption.GenerateAESKeyAndIV()
	if err != nil {
		return nil, err
	}
	inputData := &model.EncryptedService{
		ID:                 uuid.New().String(),
		Description:        derefString(encryption.EncryptField(&input.Description, aesKey, iv)),
		Amount:             derefString(encryption.EncryptField(&input.Amount, aesKey, iv)),
		Rate:               derefString(encryption.EncryptField(&input.Rate, aesKey, iv)),
		Quantity:           derefString(encryption.EncryptField(&input.Quantity, aesKey, iv)),
		EncryptedReceiptID: input.EncryptedReceiptID,
		CreatedAt:          time.Now().Format(time.RFC3339),
		AesIv:              strPtr(base64.StdEncoding.EncodeToString(iv)),
		AesKeyEncrypted:    strPtr(string(encryption.EncryptKey(r.publicKeyPEM, aesKey))),
	}
	if err := r.db.Create(inputData).Error; err != nil {
		return nil, err
	}
	return inputData, nil
}

func (r *EncryptedServiceResolver) UpdateEncryptedService(ctx context.Context, input model.UpdateEncryptedService) (*model.EncryptedService, error) {
	service := &model.EncryptedService{
		ID: input.ID,
	}

	if err := r.db.Model(service).Updates(input).Error; err != nil {
		return nil, err
	}
	newEncryptedService, err := r.GetEncryptedServiceFromDB(input.ID)
	if err != nil {
		return nil, err
	}

	return newEncryptedService, nil
}

func (r *EncryptedServiceResolver) DeleteEncryptedService(ctx context.Context, id string) (bool, error) {
	if err := r.DeleteEncryptedServiceFromDB(ctx, id); err != nil {
		return false, err
	}

	return true, nil
}
