package receiptPDFGenerator

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/loaders"
	"github.com/bishal-dd/receipt-generator-backend/graph/model"
)

func (r *ReceiptPDFGeneratorResolver) LoadServiceFromReceipt(ctx context.Context, receipt *model.Receipt) (*model.Receipt, error) {
	loaders := loaders.For(ctx)
    serviceResults, err := loaders.ServiceLoader.Load(ctx, receipt.ID)
    if err != nil {
        return nil, err
    }
        receipt.Services = serviceResults  
    
	return receipt, nil
}