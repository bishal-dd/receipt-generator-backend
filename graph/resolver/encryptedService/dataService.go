package encryptedservice

import (
	"context"
	"fmt"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/encryption"
)

func (r *EncryptedServiceResolver) GetEncryptedServiceFromDB(id string) (*model.EncryptedService, error) {
	var service *model.EncryptedService
	if err := r.db.Where("id = ?", id).First(&service).Error; err != nil {
		return nil, err
	}
	return service, nil
}

func (r *EncryptedServiceResolver) DeleteEncryptedServiceFromDB(ctx context.Context, id string) error {
	service := &model.EncryptedService{
		ID: id,
	}
	if err := r.db.Delete(service).Error; err != nil {
		return err
	}
	return nil
}

func (r *EncryptedServiceResolver) decryptService(service *model.EncryptedService) error {
	aesKey, iv, err := encryption.DecryptKeyAndIV(r.privateKeyPEM, *service.AesKeyEncrypted, *service.AesIv)
	if err != nil {
		return fmt.Errorf("decrypt AES key: %w", err)
	}
	// 4. Apply decryption
	service.Description = derefString(encryption.DecryptField(strPtr(service.Description), aesKey, iv))
	service.Amount = derefString(encryption.DecryptField(strPtr(service.Amount), aesKey, iv))
	service.Quantity = derefString(encryption.DecryptField(strPtr(service.Quantity), aesKey, iv))
	service.Rate = derefString(encryption.DecryptField(strPtr(service.Rate), aesKey, iv))

	return nil
}
