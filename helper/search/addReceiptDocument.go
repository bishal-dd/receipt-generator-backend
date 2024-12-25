package search

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/go-resty/resty/v2"
)

func AddReceiptDocument (httpClient *resty.Client, receipt model.Receipt ) error {
	typeSenseURL := os.Getenv("TYPESENSE_URL")
	if typeSenseURL == "" {
		return fmt.Errorf("TYPESENSE_URL is not set")
	}

	// Replace nil fields with empty strings
	recipientEmail := ensureString(receipt.RecipientEmail)
	recipientAddress := ensureString(receipt.RecipientAddress)
	recipientPhone := ensureString(receipt.RecipientPhone)
	paymentNote := ensureString(receipt.PaymentNote)
	subTotalAmount := ensureFloat(receipt.SubTotalAmount)
	taxAmount := ensureFloat(receipt.TaxAmount)
	year, err := getYearFromDate(receipt.Date)
	if err != nil {
		return err
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
		"recipient_email": 	recipientEmail,
		"recipient_address": 	recipientAddress,
		"recipient_phone": 	recipientPhone,
		"receipt_no": 	receipt.ReceiptNo,
		"payment_method": 	receipt.PaymentMethod,
		"payment_note": 	paymentNote,
		"sub_total_amount": 	subTotalAmount,
		"tax_amount": 	taxAmount,
		"year": year,
		"created_at": 	receipt.CreatedAt,
		"updated_at": 	receipt.UpdatedAt,
		"deleted_at": 	receipt.DeletedAt,
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

// Helper function to ensure a string is not nil
func ensureString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

// Helper function to ensure a float is not nil
func ensureFloat(value *float64) float64 {
	if value == nil {
		return 0
	}
	return *value
}

// Helper function to get the year from a date in ISO format
func getYearFromDate(date string) (int, error) {
	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return 0, fmt.Errorf("failed to parse date: %v", err)
	}
	return parsedDate.Year(), nil
}