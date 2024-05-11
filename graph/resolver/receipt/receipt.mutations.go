package receipt

import (
	"context"
	"time"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/pkg/db"
)


func (r *ReceiptResolver) CreateReceipt(ctx context.Context, input model.CreateReceipt) (*model.Receipt, error) {
    db := db.Init()
	currentTime := time.Now() 
    newReceipt := &model.Receipt{
		ID: input.ID,
        ReceiptName:    input.ReceiptName,
        RecipientName: input.RecipientName,
        RecipientPhone: input.RecipientPhone,
        Amount: input.Amount,
        TransactionNo: input.TransactionNo,
        UserID: input.UserID,
        Date: input.Date,
        TotalAmount: input.TotalAmount,
		CreatedAt: currentTime.Format("2006-01-02 15:04:05"),
    }
    if err := db.Create(newReceipt).Error; err != nil {
        return nil, err
    }

    return newReceipt, nil
}



func (r *ReceiptResolver) DeleteReceipt(ctx context.Context, id string) (bool, error) {
    db := db.Init()
    receipt := &model.Receipt{
        ID: id,
    }
    if err := db.Delete(receipt).Error; err != nil {
        return false, err
    }

    return true, nil
}