package encryptedservice

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
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
