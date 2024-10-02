package service

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
)

func (r *ServiceResolver) GetServiceFromDB(id string) (*model.Service, error) {
	var service *model.Service
	if err := r.db.Where("id = ?", id).First(&service).Error; err != nil {
		return nil, err
	}
	return service, nil
}

func (r *ServiceResolver) DeleteServiceFromDB(ctx context.Context, id string) (error) {
	service := &model.Service{
		ID: id,
	}
	if err := r.db.Delete(service).Error; err != nil {
		return  err
	}
	return  nil
} 