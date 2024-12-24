package search

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/go-resty/resty/v2"
)

func AddReceiptDocument (httpClient *resty.Client, receipt model.Receipt ) error {
	typeSenseURL := os.Getenv("TYPESENSE_URL")
	if typeSenseURL == "" {
		return fmt.Errorf("TYPESENSE_URL is not set")
	}
	resp, err := httpClient.R().
	SetHeader("Content-Type", "application/json").
	SetHeader("X-TYPESENSE-API-KEY", os.Getenv("TYPESENSE_API_KEY")).
	SetBody(map[string]interface{}{
		"id": 		receipt.ID,
		"date": 	receipt.Date,
		"total_amount": 	receipt.TotalAmount,
		"user_id": 	receipt.UserID,
		"recipient_name": 	receipt.RecipientName,
		"recipient_email": 	receipt.RecipientEmail,
		"recipient_address": 	receipt.RecipientAddress,
		"recipient_phone": 	receipt.RecipientPhone,
		"receipt_no": 	receipt.ReceiptNo,
		"payment_method": 	receipt.PaymentMethod,
		"payment_note": 	receipt.PaymentNote,
		}).
	Post(fmt.Sprintf("%s/collections/receipts/documents", typeSenseURL ) )

if err != nil {
	return  err
}
if resp.StatusCode() != http.StatusCreated {
	return fmt.Errorf("typesense returned status %d: %s", resp.StatusCode(), resp.String())
}
return  nil
}