package receipt

import (
	"fmt"
	"strconv"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/paginationUtil"
)

func convertEdges(edges []*paginationUtil.Edge[*model.Receipt]) []*model.ReceiptEdge {
	dataEdges := make([]*model.ReceiptEdge, len(edges))
	for i, edge := range edges {
		dataEdges[i] = &model.ReceiptEdge{
			Cursor: edge.Cursor,
			Node:   edge.Node,
		}
	}
	return dataEdges
}

func searchUpdateInput(input model.UpdateReceipt) map[string]interface{} {
	updateInput := make(map[string]interface{})

	// Only add non-nil fields to the updateInput map
	if input.RecipientName != nil {
		updateInput["recipient_name"] = input.RecipientName
	}
	if input.RecipientEmail != nil {
		updateInput["recipient_email"] = input.RecipientEmail
	}
	if input.RecipientAddress != nil {
		updateInput["recipient_address"] = input.RecipientAddress
	}
	if input.RecipientPhone != nil {
		updateInput["recipient_phone"] = input.RecipientPhone
	}
	if input.ReceiptNo != nil {
		updateInput["receipt_no"] = input.ReceiptNo
	}
	if input.PaymentMethod != nil {
		updateInput["payment_method"] = input.PaymentMethod
	}
	if input.PaymentNote != nil {
		updateInput["payment_note"] = input.PaymentNote
	}
	if input.UserID != nil {
		updateInput["user_id"] = input.UserID
	}
	if input.TotalAmount != nil {
		updateInput["total_amount"] = input.TotalAmount
	}

	return updateInput
}

func ParseStringToFloat64Ptr(s *string) (*float64, error) {
	if s == nil {
		return nil, nil
	}
	f, err := strconv.ParseFloat(*s, 64)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func EncryptedReceiptToReceipt(encryptedReceipt *model.EncryptedReceipt) (*model.Receipt, error) {
	totalAmount, err := ParseStringToFloat64Ptr(encryptedReceipt.TotalAmount)
	if err != nil {
		return nil, fmt.Errorf("parse total amount: %w", err)
	}

	subTotalAmount, err := ParseStringToFloat64Ptr(encryptedReceipt.SubTotalAmount)
	if err != nil {
		return nil, fmt.Errorf("parse subtotal amount: %w", err)
	}

	taxAmount, err := ParseStringToFloat64Ptr(encryptedReceipt.TaxAmount)
	if err != nil {
		return nil, fmt.Errorf("parse tax amount: %w", err)
	}
	receipt := &model.Receipt{
		ID:               encryptedReceipt.ID,
		ReceiptName:      encryptedReceipt.ReceiptName,
		RecipientName:    encryptedReceipt.RecipientName,
		RecipientPhone:   encryptedReceipt.RecipientPhone,
		RecipientEmail:   encryptedReceipt.RecipientEmail,
		RecipientAddress: encryptedReceipt.RecipientAddress,
		ReceiptNo:        encryptedReceipt.ReceiptNo,
		UserID:           encryptedReceipt.UserID,
		Date:             encryptedReceipt.Date,
		TotalAmount:      totalAmount,
		SubTotalAmount:   subTotalAmount,
		TaxAmount:        taxAmount,
		PaymentMethod:    encryptedReceipt.PaymentMethod,
		PaymentNote:      encryptedReceipt.PaymentNote,
		IsReceiptSend:    encryptedReceipt.IsReceiptSend,
		CreatedAt:        encryptedReceipt.CreatedAt,
		UpdatedAt:        encryptedReceipt.UpdatedAt,
		DeletedAt:        encryptedReceipt.DeletedAt,
		Services:         []*model.Service{},
	}

	return receipt, nil
}

func EncryptedReceiptTotalAmountToTotalAmount(encryptedReceipt *model.EncryptedReceipt) (*float64, error) {
	totalAmount, err := ParseStringToFloat64Ptr(encryptedReceipt.TotalAmount)
	if err != nil {
		return nil, fmt.Errorf("parse total amount: %w", err)
	}
	return totalAmount, nil
}
