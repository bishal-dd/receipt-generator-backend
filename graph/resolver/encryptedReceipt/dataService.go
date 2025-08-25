package encryptedReceipt

import (
	"context"
	"fmt"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/encryption"
	"github.com/bishal-dd/receipt-generator-backend/helper/stringUtil"
)

func (r *EncryptedReceiptResolver) CountTotalEncryptedReceipts() (int64, error) {
	var totalEncryptedReceipts int64
	if err := r.db.Model(&model.EncryptedReceipt{}).Count(&totalEncryptedReceipts).Error; err != nil {
		return 0, err
	}
	return totalEncryptedReceipts, nil
}

func (r *EncryptedReceiptResolver) FetchEncryptedReceiptsFromDB(ctx context.Context, offset, limit int, userId string) ([]*model.EncryptedReceipt, error) {
	var receipts []*model.EncryptedReceipt
	if err := r.db.Offset(offset).Limit(limit).Find(&receipts).Error; err != nil {
		return nil, err
	}
	// receipts, err := r.LoadServiceFromEncryptedReceipts(ctx, receipts)
	// if err != nil {
	// 	return nil, err
	// }
	return receipts, nil
}

func (r *EncryptedReceiptResolver) DeleteEncryptedReceiptFromDB(ctx context.Context, id string) error {
	receipt := &model.EncryptedReceipt{
		ID: id,
	}
	if err := r.db.Delete(receipt).Error; err != nil {
		return err
	}
	return nil
}

func (r *EncryptedReceiptResolver) GetEncryptedReceiptFromDB(ctx context.Context, id string) (*model.Receipt, error) {
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

	subTotalAmount, err := stringUtil.ParseStringToFloat64Ptr(subTotalAmountStr)
	if err != nil {
		return nil, fmt.Errorf("parse SubTotalAmount: %w", err)
	}
	taxAmount, err := stringUtil.ParseStringToFloat64Ptr(taxAmountStr)
	if err != nil {
		return nil, fmt.Errorf("parse TaxAmount: %w", err)
	}
	totalAmount, err := stringUtil.ParseStringToFloat64Ptr(totalAmountStr)
	if err != nil {
		return nil, fmt.Errorf("parse TotalAmount: %w", err)
	}

	var services []*model.Service
	for _, encryptedService := range encryptedReceipt.EncryptedServices {
		rate, err := stringUtil.ParseStringToFloat64Ptr(&encryptedService.Rate)
		if err != nil {
			return nil, fmt.Errorf("parse service price: %w", err)
		}

		quantity, err := stringUtil.ParseStringToFloat64Ptr(&encryptedService.Quantity)
		if err != nil {
			return nil, fmt.Errorf("parse service quantity: %w", err)
		}

		amount, err := stringUtil.ParseStringToFloat64Ptr(&encryptedService.Amount)
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

func (r *EncryptedReceiptResolver) decryptReceipt(receipt *model.EncryptedReceipt) error {
	aesKey, iv, err := encryption.DecryptKeyAndIV(r.privateKeyPEM, *receipt.AesKeyEncrypted, *receipt.AesIv)
	if err != nil {
		return fmt.Errorf("decrypt AES key: %w", err)
	}
	// 4. Apply decryption
	receipt.ReceiptName = encryption.DecryptField(receipt.ReceiptName, aesKey, iv)
	receipt.RecipientName = encryption.DecryptField(receipt.RecipientName, aesKey, iv)
	receipt.RecipientPhone = encryption.DecryptField(receipt.RecipientPhone, aesKey, iv)
	receipt.RecipientEmail = encryption.DecryptField(receipt.RecipientEmail, aesKey, iv)
	receipt.RecipientAddress = encryption.DecryptField(receipt.RecipientAddress, aesKey, iv)
	receipt.ReceiptNo = derefString(encryption.DecryptField(strPtr(receipt.ReceiptNo), aesKey, iv))
	receipt.PaymentMethod = derefString(encryption.DecryptField(strPtr(receipt.PaymentMethod), aesKey, iv))
	receipt.PaymentNote = encryption.DecryptField(receipt.PaymentNote, aesKey, iv)
	receipt.SubTotalAmount = encryption.DecryptField(receipt.SubTotalAmount, aesKey, iv)
	receipt.TaxAmount = encryption.DecryptField(receipt.TaxAmount, aesKey, iv)
	receipt.TotalAmount = encryption.DecryptField(receipt.TotalAmount, aesKey, iv)

	return nil
}
