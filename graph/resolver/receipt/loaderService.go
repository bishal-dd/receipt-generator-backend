package receipt

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/loaders"
	"github.com/bishal-dd/receipt-generator-backend/graph/model"
)


func (r *ReceiptResolver) LoadServiceFromReceipts(ctx context.Context, receipts []*model.Receipt) ([]*model.Receipt, error) {
	loaders := loaders.For(ctx)
	receiptIds := make([]string, len(receipts))
    for i, receipt := range receipts {
        receiptIds[i] = receipt.ID
    }

    serviceResults, err := loaders.ServiceLoader.LoadAll(ctx, receiptIds)
    if err != nil {
        return nil, err
    }
    for i, services := range serviceResults {
        receipts[i].Services = services  
    }
	return receipts, nil
}

func (r *ReceiptResolver) LoadServiceFromReceipt(ctx context.Context, receipt *model.Receipt) (*model.Receipt, error) {
	loaders := loaders.For(ctx)
    serviceResults, err := loaders.ServiceLoader.Load(ctx, receipt.ID)
    if err != nil {
        return nil, err
    }
        receipt.Services = serviceResults  
    
	return receipt, nil
}