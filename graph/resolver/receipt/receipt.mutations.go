package receipt

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper"
	"github.com/bishal-dd/receipt-generator-backend/helper/redisUtil"
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

    receiptsJSON, err := r.redis.Get(ctx, ReceiptsKey).Result()
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
    if err := redisUtil.CacheResult(r.redis, ctx, ReceiptsKey, cachedReceipts, 10); err != nil {
      return nil, err
    }

    return newReceipt, nil
}

func (r *ReceiptResolver) UpdateReceipt(ctx context.Context, input model.UpdateReceipt) (*model.Receipt, error) {
    receipt := &model.Receipt{
        ID: input.ID,
    }
    
    // Update the database entry
    if err := r.db.Model(receipt).Updates(input).Error; err != nil {
        return nil, err
    }

    // Retrieve the cached receipts
    receiptsJSON, err := r.redis.Get(ctx, ReceiptsKey).Result()
    var cachedReceipts []*model.Receipt
    
    if err != nil && err != redis.Nil { // Handle errors other than cache miss
        return nil, err
    }
    if err == redis.Nil {
        cachedReceipts = []*model.Receipt{} // Empty slice if no cache exists
    } else {
        if err := helper.Unmarshal([]byte(receiptsJSON), &cachedReceipts); err != nil {
            return nil, err
        }
    }

    // Find and update the receipt in the cache
    for i, r := range cachedReceipts {
        if r.ID == input.ID {
            cachedReceipts[i] = receipt
            break
        }
    }

    // Update the cache with the new receipts list
    if err := redisUtil.CacheResult(r.redis, ctx, ReceiptsKey, cachedReceipts, 10); err != nil {
        return nil, err
    }

    return receipt, nil
}


func (r *ReceiptResolver) DeleteReceipt(ctx context.Context, id string) (bool, error) {
    db := r.db
    receipt := &model.Receipt{
        ID: id,
    }
    cacheKey := ReceiptKey + id
    if err := db.Delete(receipt).Error; err != nil {
        return false, err
    }
    if err := r.redis.Del(ctx, cacheKey).Err(); err != nil {
        return false, err
    }
       // Get all page cache keys from the Redis set
    pageKeys, err := r.redis.SMembers(ctx, ReceiptsPageGroupKey).Result()
    if err != nil && err != redis.Nil {
        return false, err
    }
    // Delete all cached pages
    if len(pageKeys) > 0 {
        if err := r.redis.Del(ctx, pageKeys...).Err(); err != nil {
            return false, err
        }
        if err := r.redis.Del(ctx, ReceiptsPageGroupKey).Err(); err != nil {
        return false, err
    }

    }
    // Clear the set of page keys
    if err := r.redis.Del(ctx, ReceiptsPageGroupKey).Err(); err != nil {
        return false, err
    }
    return true, nil
}