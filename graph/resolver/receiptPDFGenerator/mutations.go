package receiptPDFGenerator

import (
	"context"
	"fmt"
	"log"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/contextUtil"
	"github.com/clerk/clerk-sdk-go/v2/organization"
)


func (r *ReceiptPDFGeneratorResolver) CreateReceiptPDFGenerator(ctx context.Context, input model.CreateReceiptPDFGenerator) (bool, error) {
    userId, err := contextUtil.UserIdFromContext(ctx)
    if err != nil {
        return false, err
    }
	organization, err := organization.Get(ctx, input.OrginazationID)
	if err != nil {
		return false, err
	}
	profile, err := r.GetProfileByUserID(userId)
    if err != nil {
        fmt.Println("Could not fetch profile:", err)
    }
	if err := updateProfileImages(profile, organization); err != nil {
		return false, err
	}
	totalAmount, subtotal, taxAmount := calculateTotalAmount(input, profile.Tax)
    receiptModel := inputToReceiptModel(input, userId, totalAmount, subtotal, taxAmount)
	if err := r.saveReceipt(receiptModel, input); err != nil {
		return false, err
	}
	fileName, pdfFile, err := r.generatePDF(receiptModel, profile)
	if err != nil {
        return false, err
    }
	fileURL,err := r.getFileURL(pdfFile, fileName, organization.ID, userId)
	if err != nil {
		return false, err
	}
	if receiptModel.RecipientPhone != 0 {
        err = r.sendPDFToWhatsApp(fileURL,fileName, organization.Name)
        if err != nil {
            log.Printf("Failed to send PDF to WhatsApp: %v", err)
			return false, fmt.Errorf("failed to send pdf to whatsapp: %v", err)
        }
    }
    return true, nil
}

