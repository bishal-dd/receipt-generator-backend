package receiptPDFGenerator

import (
	"fmt"
	"time"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/ids"
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
	parsedDate, err := time.Parse(time.RFC3339, receiptModel.Date)
	if err != nil {
		tx.Rollback()
		return  fmt.Errorf("invalid date format in receipt: %w", err)
	}
	receiptModel.Date = parsedDate.Format("2 January 2006")
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