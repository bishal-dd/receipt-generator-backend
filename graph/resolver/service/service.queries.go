package service

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
)

func (r *ServiceResolver) ServiceByReceiptID(ctx context.Context, receiptId string) ([]*model.Service, error) {
	var services []*model.Service
	if err := r.db.Where("receipt_id = ?", receiptId).First(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}


func (r *ServiceResolver) Service(ctx context.Context, id string) (*model.Service, error) {
	service, err := r.GetServiceFromDB(id) 
	if err != nil {
		return nil, err
	}
	return service, nil
}