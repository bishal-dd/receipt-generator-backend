package receipt

// func (r *ReceiptResolver) GetCachedReceiptPages(ctx context.Context, userId string, offset, limit int) ([]*model.Receipt, error) {
//     pageCacheKey := fmt.Sprintf("%s:%d:%d:%s", ReceiptsKey, offset, limit, userId)
//     receiptsJSON, err := r.redis.Get(ctx, pageCacheKey).Result()
//     if err == redis.Nil {
//         return nil, nil
//     } else if err != nil {
//         return nil, err
//     }

//     var receipts []*model.Receipt
//     if err := json.Unmarshal([]byte(receiptsJSON), &receipts); err != nil {
//         return nil, err
//     }
//     return receipts, nil
// }

// func (r *ReceiptResolver) GetCachedReceipt(ctx context.Context, id string) (*model.Receipt, error) {
// 	var receipt *model.Receipt
// 	cacheKey := ReceiptKey + id
// 	receiptJSON, err := r.redis.Get(ctx, cacheKey).Result()
// 	if err == redis.Nil {
// 		return nil, nil
// 	}
// 	if err := json.Unmarshal([]byte(receiptJSON), &receipt); err != nil {
// 		return nil, err
// 	}

// 	 return receipt, nil
// }