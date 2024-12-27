package receipt

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/database"
	"github.com/bishal-dd/receipt-generator-backend/helper/search"
)


func (r *ReceiptResolver) CreateReceipt(ctx context.Context, input model.CreateReceipt) (*model.Receipt, error) {
   
    inputData := database.CreateFields[model.Receipt](input);
   
    if err := r.db.Create(inputData).Error; err != nil {
        return nil, err
    }
    if err := search.AddReceiptDocument(r.httpClient, *inputData); err != nil {
        return nil, err
    }
    return inputData, nil
}

func (r *ReceiptResolver) UpdateReceipt(ctx context.Context, input model.UpdateReceipt) (*model.Receipt, error) {
    receipt := &model.Receipt{
        ID: input.ID,
    }

    updateInput := searchUpdateInput(input)
    
    if err := r.db.Model(receipt).Updates(input).Error; err != nil {
        return nil, err
    }
    if err := search.UpdateReceiptDocument(r.httpClient, updateInput, input.ID); err != nil {
        return nil, err
    }
    newReceipt, err := r.GetReceiptFromDB(ctx, input.ID)
    if err != nil {
        return nil, err
    }

    return newReceipt, nil
}


func (r *ReceiptResolver) DeleteReceipt(ctx context.Context, id string) (bool, error) {
   
    if err := r.DeleteReceiptFromDB(ctx, id); err != nil {
        return false, err
    }
    if err := search.DeleteReceiptDocument(r.httpClient, id); err != nil {
        return false, err
    }

    return true, nil
}