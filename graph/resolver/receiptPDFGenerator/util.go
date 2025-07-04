package receiptPDFGenerator

import (
	"errors"
	"time"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/cloudFront"
	"github.com/bishal-dd/receipt-generator-backend/helper/ids"
	"github.com/clerk/clerk-sdk-go/v2"
)

func calculateTotalAmount(services []*model.CreateBulkService, tax float64) (float64, float64, float64) {
	subtotal := 0.0
    for _, serviceInput := range services {
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
	profile.CompanyName = &organization.Name
	if (organization.HasImage){
		profile.LogoImage = organization.ImageURL
	}

	return nil
}

func emailInputToReceiptModel(input model.SendReceiptPDFToEmail, userId string, totalAmount, subtotal, taxAmount float64) *model.Receipt {
	receiptInput := input

	
	return &model.Receipt{
		ID:        ids.UUID(),
        UserID:    userId,
        CreatedAt: time.Now().Format(time.RFC3339),
		ReceiptName: receiptInput.ReceiptName,
		RecipientPhone: receiptInput.RecipientPhone,
		RecipientName: receiptInput.RecipientName,
		RecipientEmail: &receiptInput.RecipientEmail,
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

func whatsAppInputToReceiptModel(input model.SendReceiptPDFToWhatsApp, userId string, totalAmount, subtotal, taxAmount float64) *model.Receipt {
	receiptInput := input

	return &model.Receipt{
		ID:        ids.UUID(),
        UserID:    userId,
        CreatedAt: time.Now().Format(time.RFC3339),
		ReceiptName: receiptInput.ReceiptName,
		RecipientPhone: &receiptInput.RecipientPhone,
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


func downloadInputToReceiptModel(input model.DownloadPDF, userId string, totalAmount, subtotal, taxAmount float64) *model.Receipt {
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

func checkUserMode(mode string, useCount int) error {
	if mode == trial && useCount >= trialUseCountLimit {
		return errors.New("trial limit exceeded")
	}
	if mode == starter && useCount >= starterUseCountLimit {
		return errors.New("starter limit exceeded")
	}
	if mode == growth && useCount >= growthUseCountLimit {
		return errors.New("growth limit exceeded")
	}
	return nil
}