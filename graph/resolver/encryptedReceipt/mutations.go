package encryptedReceipt

import (
	"context"
	"crypto/aes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/google/uuid"
)

func (r *EncryptedReceiptResolver) CreateEncryptedReceipt(ctx context.Context, input model.CreateEncryptedReceipt) (*model.EncryptedReceipt, error) {
	// Step 1: Generate AES key and IV
	aesKey := make([]byte, 32) // AES-256
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(aesKey); err != nil {
		return nil, err
	}
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	// Step 2: Encrypt each field
	encryptField := func(value *string) *string {
		if value == nil {
			return nil
		}
		encrypted, _ := encryptAES(aesKey, iv, []byte(*value))
		base64Str := base64.StdEncoding.EncodeToString(encrypted)
		return &base64Str
	}

	totalAmount := ""
	if input.TotalAmount != "" {
		totalAmount = fmt.Sprintf("%s", *&input.TotalAmount)
	}

	inputData := &model.EncryptedReceipt{
		ID:               uuid.New().String(),
		UserID:           input.UserID,
		Date:             input.Date,
		IsReceiptSend:    input.IsReceiptSend,
		CreatedAt:        time.Now().Format(time.RFC3339),
		ReceiptName:      encryptField(input.ReceiptName),
		RecipientName:    encryptField(input.RecipientName),
		RecipientPhone:   encryptField(input.RecipientPhone),
		RecipientEmail:   encryptField(input.RecipientEmail),
		RecipientAddress: encryptField(input.RecipientAddress),
		ReceiptNo:        derefString(encryptField(&input.ReceiptNo)),
		PaymentMethod:    derefString(encryptField(&input.PaymentMethod)),
		PaymentNote:      encryptField(input.PaymentNote),
		TotalAmount:      encryptField(&totalAmount),
		AesIv:            strPtr(base64.StdEncoding.EncodeToString(iv)),
		AesKeyEncrypted:  strPtr(string(encryptRSAWithBase64(r.publicKeyPEM, aesKey))),
		// already base64 encoded
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
