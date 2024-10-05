package serviceLoader

import (
	"context"
	"fmt"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"gorm.io/gorm"
)


type ServiceReader struct {
	db *gorm.DB
}

func NewServiceReader(db *gorm.DB) *ServiceReader {
	return &ServiceReader{db: db}
}
func (u *ServiceReader) GetServicesByReceiptIds(ctx context.Context, receiptIds []string) ([][]*model.Service, []error) {
    if len(receiptIds) == 0 {
        return [][]*model.Service{}, nil
    }

    var services []*model.Service
    err := u.db.Where("receipt_id IN (?)", receiptIds).Find(&services).Error
    if err != nil {
        return nil, []error{fmt.Errorf("failed to fetch services: %w", err)}
    }

    serviceMap := make(map[string][]*model.Service)
    for _, service := range services {
        serviceMap[service.ReceiptID] = append(serviceMap[service.ReceiptID], service)
    }

    result := make([][]*model.Service, len(receiptIds))
    for i, userID := range receiptIds {
        result[i] = serviceMap[userID] 
    }

    return result, nil
}
