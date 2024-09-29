package receipt

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper"
	"github.com/bishal-dd/receipt-generator-backend/helper/contextUtil"
	"github.com/bishal-dd/receipt-generator-backend/helper/redisUtil"
)


func (r *ReceiptResolver) CreateReceipt(ctx context.Context, input model.CreateReceipt) (*model.Receipt, error) {
    newReceipt := &model.Receipt{
        ID: helper.UUID(),
        ReceiptName:    input.ReceiptName,
        RecipientName: input.RecipientName,
        RecipientPhone: input.RecipientPhone,
        Amount: input.Amount,
        TransactionNo: input.TransactionNo,
        UserID: input.UserID,
        Date: input.Date,
        TotalAmount: input.TotalAmount,
		CreatedAt: helper.CurrentTime(),
    }
    if err := r.db.Create(newReceipt).Error; err != nil {
        return nil, err
    }

    return newReceipt, nil
}

func (r *ReceiptResolver) UpdateReceipt(ctx context.Context, input model.UpdateReceipt) (*model.Receipt, error) {
    receipt := &model.Receipt{
        ID: input.ID,
    }
    
    if err := r.db.Model(receipt).Updates(input).Error; err != nil {
        return nil, err
    }
    if err := redisUtil.DeleteCacheItem(r.redis, ctx, ReceiptKey, input.ID); err != nil {
        return nil, err
    }
    newReceipt, err := r.GetReceiptFromDB(input.ID)
    if err != nil {
        return nil, err
    }

    return newReceipt, nil
}


func (r *ReceiptResolver) DeleteReceipt(ctx context.Context, id string) (bool, error) {
    userId, err := contextUtil.UserIdFromContext(ctx)
    if err != nil {
        return false, err
    }
    if err := r.DeleteReceiptFromDB(ctx, id); err != nil {
        return false, err
    }
    if err := redisUtil.DeleteCacheItem(r.redis, ctx, ReceiptKey, id); err != nil {
        return false, err
    }
    if err := redisUtil.DeleteCachePages(r.redis, ctx, ReceiptsPageGroupKey, userId); err != nil {
        return false, err
    }

    return true, nil
}