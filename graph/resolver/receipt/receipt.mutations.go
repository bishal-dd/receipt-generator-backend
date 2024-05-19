package receipt

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper"
	"github.com/redis/go-redis/v9"
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

    receiptsJSON, err := r.redis.Get(ctx, Key).Result()
    var cachedReceipts []*model.Receipt
  
    if err != nil && err != redis.Nil {  // Handle errors other than cache miss
      return nil, err
    }
    if err == redis.Nil {
      cachedReceipts = []*model.Receipt{}  // Empty slice if no cache exists
    } else {
      if err := helper.Unmarshal([]byte(receiptsJSON), &cachedReceipts); err != nil {
        return nil, err
      }
    }
    cachedReceipts = append(cachedReceipts, newReceipt)
    if err := helper.CacheResult(r.redis, ctx, Key, cachedReceipts, 10); err != nil {
      return nil, err
    }

    return newReceipt, nil
}



func (r *ReceiptResolver) DeleteReceipt(ctx context.Context, id string) (bool, error) {
    db := r.db
    receipt := &model.Receipt{
        ID: id,
    }
    if err := db.Delete(receipt).Error; err != nil {
        return false, err
    }

    return true, nil
}