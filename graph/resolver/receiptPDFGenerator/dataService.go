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

func (r *ReceiptPDFGeneratorResolver) GetEncryptedReceiptFromDB(ctx context.Context, id string) (*model.Receipt, error) {
	var encryptedReceipt *model.EncryptedReceipt
	if err := r.db.Where("id = ?", id).First(&encryptedReceipt).Error; err != nil {
		fmt.Print(err)
		return nil, err
	}

	encryptedReceipt, err := r.LoadServiceFromEncryptedReceipt(ctx, encryptedReceipt)
	if err != nil {
		return nil, err
	}
	aesKey, iv, err := encryption.DecryptKeyAndIV(r.privateKeyPEM, *encryptedReceipt.AesKeyEncrypted, *encryptedReceipt.AesIv)
	if err != nil {
		return nil, fmt.Errorf("decrypt AES key: %w", err)
	}

	subTotalAmountStr := encryption.DecryptField(encryptedReceipt.SubTotalAmount, aesKey, iv)
	taxAmountStr := encryption.DecryptField(encryptedReceipt.TaxAmount, aesKey, iv)
	totalAmountStr := encryption.DecryptField(encryptedReceipt.TotalAmount, aesKey, iv)

	subTotalAmount, err := ParseStringToFloat64Ptr(subTotalAmountStr)
	if err != nil {
		return nil, fmt.Errorf("parse SubTotalAmount: %w", err)
	}
	taxAmount, err := ParseStringToFloat64Ptr(taxAmountStr)
	if err != nil {
		return nil, fmt.Errorf("parse TaxAmount: %w", err)
	}
	totalAmount, err := ParseStringToFloat64Ptr(totalAmountStr)
	if err != nil {
		return nil, fmt.Errorf("parse TotalAmount: %w", err)
	}

	var services []*model.Service
	for _, encryptedService := range encryptedReceipt.EncryptedServices {
		rate, err := ParseStringToFloat64Ptr(&encryptedService.Rate)
		if err != nil {
			return nil, fmt.Errorf("parse service price: %w", err)
		}

		quantity, err := ParseStringToFloat64Ptr(&encryptedService.Quantity)
		if err != nil {
			return nil, fmt.Errorf("parse service quantity: %w", err)
		}

		amount, err := ParseStringToFloat64Ptr(&encryptedService.Amount)
		if err != nil {
			return nil, fmt.Errorf("parse service amount: %w", err)
		}

		services = append(services, &model.Service{
			ID:          encryptedService.ID,
			Description: encryptedService.Description,
			Rate:        *rate,
			Quantity:    int(*quantity),
			Amount:      *amount,
			ReceiptID:   encryptedService.EncryptedReceiptID,
			CreatedAt:   encryptedService.CreatedAt,
			DeletedAt:   encryptedService.DeletedAt,
			UpdatedAt:   encryptedService.UpdatedAt,
		})
	}

	receipt := &model.Receipt{
		ID:               encryptedReceipt.ID,
		RecipientPhone:   encryption.DecryptField(encryptedReceipt.RecipientPhone, aesKey, iv),
		RecipientEmail:   encryption.DecryptField(encryptedReceipt.RecipientEmail, aesKey, iv),
		RecipientAddress: encryption.DecryptField(encryptedReceipt.RecipientAddress, aesKey, iv),
		ReceiptNo:        stringUtil.DerefString(encryption.DecryptField(stringUtil.StrPtr(encryptedReceipt.ReceiptNo), aesKey, iv)),
		PaymentMethod:    stringUtil.DerefString(encryption.DecryptField(stringUtil.StrPtr(encryptedReceipt.PaymentMethod), aesKey, iv)),
		PaymentNote:      encryption.DecryptField(encryptedReceipt.PaymentNote, aesKey, iv),
		SubTotalAmount:   subTotalAmount,
		Date:             encryptedReceipt.Date,
		Services:         services,
		TaxAmount:        taxAmount,
		TotalAmount:      totalAmount,
		CreatedAt:        encryptedReceipt.CreatedAt,
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

func (r *ReceiptPDFGeneratorResolver) saveReceiptFile(encryptedReceiptModel *model.EncryptedReceipt, tx *gorm.DB) (*model.ReceiptFile, error) {
	receiptFile := &model.ReceiptFile{
		ID:                 ids.UUID(),
		ReceiptNo:          encryptedReceiptModel.ReceiptNo,
		EncryptedReceiptID: encryptedReceiptModel.ID,
		IssuedAt:           time.Now().Format(time.RFC3339),
		CreatedAt:          time.Now().Format(time.RFC3339),
		UpdatedAt:          time.Now().Format(time.RFC3339),
	}

	if err := tx.Table("receipt_files").Create(receiptFile).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return receiptFile, nil
}

func (r *ReceiptPDFGeneratorResolver) saveReceiptFileFromReceiptModel(receiptModel *model.Receipt) (*model.ReceiptFile, error) {
	receiptFile := &model.ReceiptFile{
		ID:                 ids.UUID(),
		ReceiptNo:          receiptModel.ReceiptNo,
		EncryptedReceiptID: receiptModel.ID,
		IssuedAt:           time.Now().Format(time.RFC3339),
		CreatedAt:          time.Now().Format(time.RFC3339),
		UpdatedAt:          time.Now().Format(time.RFC3339),
	}

	if err := r.db.Table("receipt_files").Create(receiptFile).Error; err != nil {
		return nil, err
	}

	return receiptFile, nil
}
