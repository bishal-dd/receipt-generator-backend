package encryptedReceipt

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
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
	// 1. Decrypt AES key
	aesKey, err := decryptRSAWithBase64([]byte(r.privateKeyPEM), []byte(*receipt.AesKeyEncrypted))
	if err != nil {
		return fmt.Errorf("decrypt AES key: %w", err)
	}

	// 2. Decode IV
	iv, err := base64.StdEncoding.DecodeString(*receipt.AesIv)
	if err != nil {
		return fmt.Errorf("decode IV: %w", err)
	}

	// 3. Decrypt helper
	decryptField := func(encryptedBase64 *string) *string {
		if encryptedBase64 == nil {
			return nil
		}
		encryptedBytes, err := base64.StdEncoding.DecodeString(*encryptedBase64)
		if err != nil {
			return nil
		}
		decrypted, err := decryptAES(aesKey, iv, encryptedBytes)
		if err != nil {
			return nil
		}
		str := string(decrypted)
		return &str
	}

	// 4. Apply decryption
	receipt.ReceiptName = decryptField(receipt.ReceiptName)
	receipt.RecipientName = decryptField(receipt.RecipientName)
	receipt.RecipientPhone = decryptField(receipt.RecipientPhone)
	receipt.RecipientEmail = decryptField(receipt.RecipientEmail)
	receipt.RecipientAddress = decryptField(receipt.RecipientAddress)
	receipt.ReceiptNo = derefString(decryptField(strPtr(receipt.ReceiptNo)))
	receipt.PaymentMethod = derefString(decryptField(strPtr(receipt.PaymentMethod)))
	receipt.PaymentNote = decryptField(receipt.PaymentNote)
	receipt.TotalAmount = decryptField(receipt.TotalAmount)

	return nil
}
