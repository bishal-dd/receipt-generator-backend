package receiptPDFGenerator

import (
	"time"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/cloudFront"
	"github.com/bishal-dd/receipt-generator-backend/helper/ids"
	"github.com/clerk/clerk-sdk-go/v2"
)

func calculateTotalAmount(input model.CreateReceiptPDFGenerator, tax float64) (float64, float64, float64) {
	subtotal := 0.0
    for _, serviceInput := range input.Services {
		subtotal += serviceInput.Amount
    }

	taxRate := tax / 100
	taxAmount := subtotal * taxRate
	totalAmount := subtotal + taxAmount

	return totalAmount, subtotal, taxAmount
}

func updateProfileImages(profile *model.Profile, organization *clerk.Organization ) error {
	if profile.SignatureImage != nil && *profile.SignatureImage != "" {
		signedURL, err := cloudFront.GetCloudFrontURL(*profile.SignatureImage)
		if err != nil {
			return  err
		}
		profile.SignatureImage = &signedURL
	}
	profile.CompanyName = organization.Name
	if (organization.HasImage){
		profile.LogoImage = organization.ImageURL
	}

	return nil
}

func inputToReceiptModel(input model.CreateReceiptPDFGenerator, userId string, totalAmount, subtotal, taxAmount float64) *model.Receipt {
	receiptInput := input

	return &model.Receipt{
		ID:        ids.UUID(),
        UserID:    userId,
        CreatedAt: time.Now().Format(time.RFC3339),
		ReceiptName: receiptInput.ReceiptName,
		RecipientPhone: receiptInput.RecipientPhone,
		RecipientName: receiptInput.RecipientName,
		RecipientEmail: receiptInput.RecipientEmail,
		RecipientAddress: receiptInput.RecipientAddress,
		ReceiptNo: receiptInput.ReceiptNo,
		Date: receiptInput.Date,
		PaymentMethod: receiptInput.PaymentMethod,
		PaymentNote: receiptInput.PaymentNote,
		TotalAmount: &totalAmount,
		SubTotalAmount: &subtotal,
		TaxAmount: &taxAmount,
		Services:  make([]*model.Service, 0),
	}
}

