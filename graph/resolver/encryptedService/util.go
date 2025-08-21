package encryptedservice

import (
	"fmt"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/encryption"
)

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
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
