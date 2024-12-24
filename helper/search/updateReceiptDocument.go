package search

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/go-resty/resty/v2"
)

func UpdateReceiptDocument (httpClient *resty.Client, input model.UpdateReceipt ) error {
	typeSenseURL := os.Getenv("TYPESENSE_URL")
	if typeSenseURL == "" {
		return fmt.Errorf("TYPESENSE_URL is not set")
	}
	resp, err := httpClient.R().
	SetHeader("Content-Type", "application/json").
	SetHeader("X-TYPESENSE-API-KEY", os.Getenv("TYPESENSE_API_KEY")).
	SetBody(map[string]interface{}{
		"recipient_name": 	input.RecipientName,
		"recipient_email": 	input.RecipientEmail,
		"recipient_address": 	input.RecipientAddress,
		"recipient_phone": 	input.RecipientPhone,
		"receipt_no": 	input.ReceiptNo,
		"payment_method": 	input.PaymentMethod,
		"payment_note": 	input.PaymentNote,
		"user_id": 	input.UserID,
		"total_amount": 	input.TotalAmount,
	}).
	Patch(fmt.Sprintf("%s/collections/receipts/documents/%s", typeSenseURL, input.ID ) )
	if err != nil {
		return  err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("typesense returned status %d: %s", resp.StatusCode(), resp.String())
	}
	return  nil
}