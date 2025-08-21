package encryptedReceipt

import (
	"context"
	"fmt"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/encryption"
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

func (r *EncryptedReceiptResolver) GetEncryptedReceiptFromDB(ctx context.Context, id string) (*model.EncryptedReceipt, error) {
	var receipt *model.EncryptedReceipt
	if err := r.db.Where("id = ?", id).First(&receipt).Error; err != nil {
		return nil, err
	}

	// receipt, err := r.LoadServiceFromEncryptedReceipt(ctx, receipt)
	// if err != nil {
	// 	return nil, err
	// }
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
