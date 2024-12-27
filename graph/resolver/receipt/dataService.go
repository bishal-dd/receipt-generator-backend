package receipt

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
)

func (r *ReceiptResolver) CountTotalReceipts() (int64, error) {
    var totalReceipts int64
    if err := r.db.Model(&model.Receipt{}).Count(&totalReceipts).Error; err != nil {
        return 0, err
    }
    return totalReceipts, nil
}

func (r *ReceiptResolver) FetchReceiptsFromDB(ctx context.Context, offset, limit int, userId string) ([]*model.Receipt, error) {
    var receipts []*model.Receipt
    if err := r.db.Offset(offset).Limit(limit).Find(&receipts).Error; err != nil {
        return nil, err
    }
	receipts, err := r.LoadServiceFromReceipts(ctx, receipts)
	if err != nil {
		return nil, err
	}
    return receipts, nil
}

func (r *ReceiptResolver) DeleteReceiptFromDB(ctx context.Context, id string) (error) {
	receipt := &model.Receipt{
		ID: id,
	}
	if err := r.db.Delete(receipt).Error; err != nil {
		return  err
	}
	return  nil
}

func (r *ReceiptResolver) GetReceiptFromDB(ctx context.Context, id string) (*model.Receipt, error) {
	var receipt *model.Receipt
	if err := r.db.Where("id = ?", id).First(&receipt).Error; err != nil {
		return nil, err
	}

	receipt, err := r.LoadServiceFromReceipt(ctx, receipt)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}