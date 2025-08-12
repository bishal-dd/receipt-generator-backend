package encryptedReceipt

// func (r *EncryptedReceiptResolver) LoadServiceFromEncryptedReceipts(ctx context.Context, receipts []*model.EncryptedReceipt) ([]*model.EncryptedReceipt, error) {
// 	loaders := loaders.For(ctx)
// 	receiptIds := make([]string, len(receipts))
// 	for i, receipt := range receipts {
// 		receiptIds[i] = receipt.ID
// 	}

// 	serviceResults, err := loaders.ServiceLoader.LoadAll(ctx, receiptIds)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for i, services := range serviceResults {
// 		receipts[i].Services = services
// 	}
// 	return receipts, nil
// }

// func (r *EncryptedReceiptResolver) LoadServiceFromEncryptedReceipt(ctx context.Context, receipt *model.EncryptedReceipt) (*model.EncryptedReceipt, error) {
// 	loaders := loaders.For(ctx)
// 	serviceResults, err := loaders.ServiceLoader.Load(ctx, receipt.ID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	receipt.Services = serviceResults

// 	return receipt, nil
// }
