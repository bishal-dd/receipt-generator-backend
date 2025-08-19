package receiptPDFGenerator

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/encryption"
	"github.com/bishal-dd/receipt-generator-backend/helper/ids"
	"github.com/bishal-dd/receipt-generator-backend/helper/stringUtil"
	"gorm.io/gorm"
)

func (r *ReceiptPDFGeneratorResolver) GetProfileByUserID(userId string) (*model.Profile, error) {
	var profile *model.Profile
	if err := r.db.Where("user_id = ?", userId).First(&profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

func (r *ReceiptPDFGeneratorResolver) saveReceipt(receiptModel *model.Receipt, services []*model.CreateBulkService, tx *gorm.DB) error {
	if tx.Error != nil {
		return tx.Error
	}
	// Create Receipt
	if err := tx.Create(receiptModel).Error; err != nil {
		tx.Rollback()
		return err
	}
	// if err := search.AddReceiptDocument(r.httpClient, *receiptModel); err != nil {
	//     tx.Rollback()
	//     return  err
	// }

	// Create Services
	for _, serviceInput := range services {
		serviceModel := &model.Service{
			ID:          ids.UUID(),
			ReceiptID:   receiptModel.ID,
			CreatedAt:   time.Now().Format(time.RFC3339),
			Description: serviceInput.Description,
			Rate:        serviceInput.Rate,
			Quantity:    serviceInput.Quantity,
			Amount:      serviceInput.Amount,
		}
		if err := tx.Create(serviceModel).Error; err != nil {
			tx.Rollback()
			return err
		}
		receiptModel.Services = append(receiptModel.Services, serviceModel)
	}

	return nil
}

func (r *ReceiptPDFGeneratorResolver) saveEncryptedReceipt(encryptedReceiptModel *model.EncryptedReceipt, services []*model.CreateBulkService, tx *gorm.DB) error {
	if tx.Error != nil {
		return tx.Error
	}
	// Create Receipt
	if err := tx.Create(encryptedReceiptModel).Error; err != nil {
		log.Println("Error saving encrypted service:", err)

		tx.Rollback()
		return err
	}
	// if err := search.AddReceiptDocument(r.httpClient, *receiptModel); err != nil {
	//     tx.Rollback()
	//     return  err
	// }

	// Create Services
	for _, serviceInput := range services {
		aesKey, iv, err := encryption.GenerateAESKeyAndIV()
		if err != nil {
			log.Println("Error saving encrypted service:", err)
			return err
		}
		encryptedServiceModel := &model.EncryptedService{
			ID:                 ids.UUID(),
			EncryptedReceiptID: encryptedReceiptModel.ID,
			CreatedAt:          time.Now().Format(time.RFC3339),
			Description:        stringUtil.DerefString(encryption.EncryptField(&serviceInput.Description, aesKey, iv)),
			Rate:               stringUtil.DerefString(encryption.EncryptField(stringUtil.StrPtr(fmt.Sprintf("%.2f", serviceInput.Rate)), aesKey, iv)),
			Quantity:           stringUtil.DerefString(encryption.EncryptField(stringUtil.StrPtr(fmt.Sprintf("%d", serviceInput.Quantity)), aesKey, iv)),
			Amount:             stringUtil.DerefString(encryption.EncryptField(stringUtil.StrPtr(fmt.Sprintf("%.2f", serviceInput.Amount)), aesKey, iv)),
			AesIv:              stringUtil.StrPtr(base64.StdEncoding.EncodeToString(iv)),
			AesKeyEncrypted:    stringUtil.StrPtr(string(encryption.EncryptKey(r.publicKeyPEM, aesKey))),
		}
		if err := tx.Create(encryptedServiceModel).Error; err != nil {
			log.Println("Error saving encrypted service:", err)
			tx.Rollback()
			return err
		}
		encryptedReceiptModel.EncryptedServices = append(encryptedReceiptModel.EncryptedServices, encryptedServiceModel)
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

func (r *ReceiptPDFGeneratorResolver) MinusProductQuantity(services []*model.CreateBulkService, tx *gorm.DB) error {
	for _, service := range services {
		var product model.Product
		if err := tx.First(&product, "id = ?", service.ID).Error; err != nil {
			tx.Rollback()
			return err
		}

		if *product.Quantity < service.Quantity {
			tx.Rollback()
			return fmt.Errorf("not enough stock for product ID %s", product.ID)
		}

		*product.Quantity -= service.Quantity

		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return nil
}
