package service

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/database"
)

func (r *ServiceResolver) CreateService(ctx context.Context, input model.CreateService) (*model.Service, error) {
	inputData := database.CreateFields[model.Service](input);

	if err := r.db.Create(inputData).Error; err != nil {
		return nil, err
	}
	return inputData, nil
}

func (r *ServiceResolver) UpdateService(ctx context.Context, input model.UpdateService) (*model.Service, error) {	
	service := &model.Service{
		ID: input.ID,
	}
	
	if err := r.db.Model(service).Updates(input).Error; err != nil {
		return nil, err
	}
	newService, err := r.GetServiceFromDB(input.ID)
	if err != nil {
		return nil, err
	}

	return newService, nil
}

func (r *ServiceResolver) DeleteService(ctx context.Context, id string) (bool, error) {
	if err := r.DeleteServiceFromDB(ctx, id); err != nil {
		return false, err
	}
	
	return true, nil
}