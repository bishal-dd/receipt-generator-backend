package receipt

import (
	"context"
	"fmt"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/encryption"
	"github.com/bishal-dd/receipt-generator-backend/helper/stringUtil"
)

func (r *ReceiptResolver) CountTotalReceipts() (int64, error) {
	var totalReceipts int64
	if err := r.db.Model(&model.Receipt{}).Count(&totalReceipts).Error; err != nil {
		return 0, err
	}
	return totalReceipts, nil
}

func (r *ReceiptResolver) FetchReceiptsFromDB(ctx context.Context, offset, limit int, userId string) ([]*model.Receipt, error) {
	var encryptedReceipts []*model.EncryptedReceipt
	var receipts []*model.Receipt
	if err := r.db.Offset(offset).Limit(limit).Find(&encryptedReceipts).Error; err != nil {
		return nil, err
	}
	for _, enc := range encryptedReceipts {
		err := r.decryptReceipt(enc)
		if err != nil {
			// Optionally skip or return error
			return nil, fmt.Errorf("failed to decrypt receipt %s: %w", enc.ID, err)
		}
		receipt, err := EncryptedReceiptToReceipt(enc)
		if err != nil {
			return nil, err
		}

		receipts = append(receipts, receipt)
	}

	receipts, err := r.LoadServiceFromReceipts(ctx, receipts)
	if err != nil {
		return nil, err
	}

	return receipts, nil
}

func (r *ReceiptResolver) decryptReceipt(receipt *model.EncryptedReceipt) error {
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
	receipt.PaymentMethod = stringUtil.DerefString(encryption.DecryptField(stringUtil.StrPtr(receipt.PaymentMethod), aesKey, iv))
	receipt.PaymentNote = encryption.DecryptField(receipt.PaymentNote, aesKey, iv)
	receipt.SubTotalAmount = encryption.DecryptField(receipt.SubTotalAmount, aesKey, iv)
	receipt.TaxAmount = encryption.DecryptField(receipt.TaxAmount, aesKey, iv)
	receipt.TotalAmount = encryption.DecryptField(receipt.TotalAmount, aesKey, iv)

	return nil
}

func (r *ReceiptResolver) DeleteReceiptFromDB(ctx context.Context, id string) error {
	receipt := &model.Receipt{
		ID: id,
	}
	if err := r.db.Delete(receipt).Error; err != nil {
		return err
	}
	return nil
}

func (r *ReceiptResolver) GetReceiptFromDB(ctx context.Context, id string) (*model.Receipt, error) {
	var encryptedReceipt model.EncryptedReceipt

	// Step 1: Fetch Encrypted Receipt from DB
	if err := r.db.Where("id = ?", id).First(&encryptedReceipt).Error; err != nil {
		return nil, err
	}

	// Step 2: Decrypt the receipt
	if err := r.decryptReceipt(&encryptedReceipt); err != nil {
		return nil, fmt.Errorf("failed to decrypt receipt %s: %w", encryptedReceipt.ID, err)
	}

	// Step 3: Convert EncryptedReceipt → Receipt

	receipt, err := EncryptedReceiptToReceipt(&encryptedReceipt)
	if err != nil {
		return nil, fmt.Errorf("failed to convert encrypted receipt %s to receipt: %w", encryptedReceipt.ID, err)
	}
	// Step 4: Load services (encrypted → decrypted → parsed)
	receipt, err = r.LoadServiceFromReceipt(ctx, receipt)
	if err != nil {
		return nil, fmt.Errorf("failed to load services for receipt %s: %w", receipt.ID, err)
	}

	return receipt, nil
}

func (r *ReceiptResolver) decryptService(service *model.EncryptedService) error {
	aesKey, iv, err := encryption.DecryptKeyAndIV(r.privateKeyPEM, *service.AesKeyEncrypted, *service.AesIv)
	if err != nil {
		return fmt.Errorf("decrypt AES key: %w", err)
	}
	// 4. Apply decryption
	service.Description = stringUtil.DerefString(encryption.DecryptField(stringUtil.StrPtr(service.Description), aesKey, iv))
	service.Amount = stringUtil.DerefString(encryption.DecryptField(stringUtil.StrPtr(service.Amount), aesKey, iv))
	service.Quantity = stringUtil.DerefString(encryption.DecryptField(stringUtil.StrPtr(service.Quantity), aesKey, iv))
	service.Rate = stringUtil.DerefString(encryption.DecryptField(stringUtil.StrPtr(service.Rate), aesKey, iv))

	return nil
}
