package receipt

import (
	"context"
	"fmt"

	"github.com/bishal-dd/receipt-generator-backend/graph/loaders"
	"github.com/bishal-dd/receipt-generator-backend/graph/model"
)

func (r *ReceiptResolver) LoadServiceFromReceipts(ctx context.Context, receipts []*model.Receipt) ([]*model.Receipt, error) {
	loaders := loaders.For(ctx)
	receiptIds := make([]string, len(receipts))
	for i, receipt := range receipts {
		receiptIds[i] = receipt.ID
	}

	// serviceResults is of type [][]*model.EncryptedService
	serviceResults, err := loaders.EncryptedServiceLoader.LoadAll(ctx, receiptIds)
	if err != nil {
		return nil, err
	}

	for _, encryptedServices := range serviceResults {
		for _, encryptedService := range encryptedServices {
			err := r.decryptService(encryptedService)
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt service %s: %w", encryptedService.ID, err)
			}
		}
	}

	for i, encryptedServices := range serviceResults {
		var services []*model.Service
		for _, encryptedService := range encryptedServices {
			rate, err := ParseStringToFloat64Ptr(&encryptedService.Rate)
			if err != nil {
				return nil, fmt.Errorf("parse service price: %w", err)
			}

			quantity, err := ParseStringToFloat64Ptr(&encryptedService.Quantity)
			if err != nil {
				return nil, fmt.Errorf("parse service quantity: %w", err)
			}

			amount, err := ParseStringToFloat64Ptr(&encryptedService.Amount)
			if err != nil {
				return nil, fmt.Errorf("parse service amount: %w", err)
			}

			services = append(services, &model.Service{
				ID:          encryptedService.ID,
				Description: encryptedService.Description,
				Rate:        *rate,
				Quantity:    int(*quantity),
				Amount:      *amount,
				ReceiptID:   encryptedService.EncryptedReceiptID,
				CreatedAt:   encryptedService.CreatedAt,
				DeletedAt:   encryptedService.DeletedAt,
				UpdatedAt:   encryptedService.UpdatedAt,
			})
		}
		// Attach services to the corresponding receipt
		receipts[i].Services = services
	}

	return receipts, nil
}

func (r *ReceiptResolver) LoadServiceFromReceipt(ctx context.Context, receipt *model.Receipt) (*model.Receipt, error) {
	loaders := loaders.For(ctx)
	serviceResults, err := loaders.EncryptedServiceLoader.Load(ctx, receipt.ID)
	if err != nil {
		return nil, err
	}
	var services []*model.Service
	for _, service := range serviceResults {
		err := r.decryptService(service)
		if err != nil {
			// Optionally skip or return error
			return nil, fmt.Errorf("failed to decrypt receipt %s: %w", service.ID, err)
		}

		rate, err := ParseStringToFloat64Ptr(&service.Rate)
		if err != nil {
			return nil, fmt.Errorf("parse service price: %w", err)
		}

		quantity, err := ParseStringToFloat64Ptr(&service.Quantity)
		if err != nil {
			return nil, fmt.Errorf("parse service quantity: %w", err)
		}

		amount, err := ParseStringToFloat64Ptr(&service.Amount)
		if err != nil {
			return nil, fmt.Errorf("parse service amount: %w", err)
		}

		services = append(services, &model.Service{
			ID:          service.ID,
			Description: service.Description,
			Rate:        *rate,
			Quantity:    int(*quantity),
			Amount:      *amount,
			ReceiptID:   service.EncryptedReceiptID,
			CreatedAt:   service.CreatedAt,
			DeletedAt:   service.DeletedAt,
			UpdatedAt:   service.UpdatedAt,
		})
	}

	receipt.Services = services

	return receipt, nil
}
