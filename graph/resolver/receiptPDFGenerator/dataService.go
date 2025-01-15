package receiptPDFGenerator

import (
	"context"
	"fmt"
	"time"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/ids"
	"github.com/bishal-dd/receipt-generator-backend/helper/search"
)

func (r *ReceiptPDFGeneratorResolver) GetProfileByUserID(userId string) (*model.Profile, error) {
	var profile *model.Profile
	if err := r.db.Where("user_id = ?", userId).First(&profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

func (r *ReceiptPDFGeneratorResolver) saveReceipt(receiptModel *model.Receipt, services []*model.CreateBulkService)error{
	tx := r.db.Begin()
    if tx.Error != nil {
        return tx.Error
    }
    // Create Receipt
    if err := tx.Create(receiptModel).Error; err != nil {
        tx.Rollback()
        return  err
    }
    if err := search.AddReceiptDocument(r.httpClient, *receiptModel); err != nil {
        fmt.Println("Error adding receipt document to search: ", err)
        tx.Rollback()
        return  err
    }
	
    // Create Services
    for _, serviceInput := range services {
        serviceModel := &model.Service{
            ID:         ids.UUID(),
            ReceiptID:  receiptModel.ID,
            CreatedAt:  time.Now().Format(time.RFC3339),
            Description: serviceInput.Description,
            Rate:      serviceInput.Rate,
            Quantity:   serviceInput.Quantity,
			Amount:    serviceInput.Amount,
        }
        if err := tx.Create(serviceModel).Error; err != nil {
            tx.Rollback()
            return err
        }
		receiptModel.Services = append(receiptModel.Services, serviceModel)
    }
    if err := tx.Commit().Error; err != nil {
        return  err
    }

	return nil
}

func (r *ReceiptPDFGeneratorResolver) GetReceiptFromDB(ctx context.Context, id string) (*model.Receipt, error) {
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


func (r *ReceiptPDFGeneratorResolver) GetUserFromDB(ctx context.Context, id string) (*model.User, error) {
    var user *model.User
    if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
        return nil, err
    }
    return user, nil
}