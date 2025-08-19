package encryptedReceipt

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/encryption"
	"github.com/bishal-dd/receipt-generator-backend/helper/ids"
)

func (r *EncryptedReceiptResolver) CreateEncryptedReceipt(ctx context.Context, input model.CreateEncryptedReceipt) (*model.EncryptedReceipt, error) {
	aesKey, iv, err := encryption.GenerateAESKeyAndIV()
	if err != nil {
		return nil, err
	}

	totalAmount := ""
	if input.TotalAmount != "" {
		totalAmount = fmt.Sprintf("%s", input.TotalAmount)
	}

	inputData := &model.EncryptedReceipt{
		ID:               ids.UUID(),
		UserID:           input.UserID,
		Date:             input.Date,
		IsReceiptSend:    input.IsReceiptSend,
		CreatedAt:        time.Now().Format(time.RFC3339),
		ReceiptName:      encryption.EncryptField(input.ReceiptName, aesKey, iv),
		RecipientName:    encryption.EncryptField(input.RecipientName, aesKey, iv),
		RecipientPhone:   encryption.EncryptField(input.RecipientPhone, aesKey, iv),
		RecipientEmail:   encryption.EncryptField(input.RecipientEmail, aesKey, iv),
		RecipientAddress: encryption.EncryptField(input.RecipientAddress, aesKey, iv),
		ReceiptNo:        derefString(encryption.EncryptField(&input.ReceiptNo, aesKey, iv)),
		PaymentMethod:    derefString(encryption.EncryptField(&input.PaymentMethod, aesKey, iv)),
		PaymentNote:      encryption.EncryptField(input.PaymentNote, aesKey, iv),
		TotalAmount:      encryption.EncryptField(&totalAmount, aesKey, iv),
		AesIv:            strPtr(base64.StdEncoding.EncodeToString(iv)),
		AesKeyEncrypted:  strPtr(string(encryption.EncryptKey(r.publicKeyPEM, aesKey))),
	}

	// Step 3: Save to DB
	if err := r.db.Create(inputData).Error; err != nil {
		return nil, err
	}

	return inputData, nil
}

func (r *EncryptedReceiptResolver) UpdateEncryptedReceipt(ctx context.Context, input model.UpdateEncryptedReceipt) (*model.EncryptedReceipt, error) {
	receipt := &model.EncryptedReceipt{
		ID: input.ID,
	}

	if err := r.db.Model(receipt).Updates(input).Error; err != nil {
		return nil, err
	}

	newEncryptedReceipt, err := r.GetEncryptedReceiptFromDB(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return newEncryptedReceipt, nil
}

func (r *EncryptedReceiptResolver) DeleteEncryptedReceipt(ctx context.Context, id string) (bool, error) {

	if err := r.DeleteEncryptedReceiptFromDB(ctx, id); err != nil {
		return false, err
	}
	// if err := search.DeleteEncryptedReceiptDocument(r.httpClient, id); err != nil {
	//     return false, err
	// }

	return true, nil
}
