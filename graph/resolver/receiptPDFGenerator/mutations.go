package receiptPDFGenerator

import (
	"context"
	"fmt"
	"log"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/contextUtil"
	"github.com/bishal-dd/receipt-generator-backend/helper/emails"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/organization"
	"golang.org/x/sync/errgroup"
)


func (r *ReceiptPDFGeneratorResolver) SendReceiptPDFToWhatsApp(ctx context.Context, input model.CreateReceiptPDFGenerator) (bool, error) {
    // Early error checking
    userId, err := contextUtil.UserIdFromContext(ctx)
    if err != nil {
        return false, err
    }

    // Use errgroup for parallel processing
    var g errgroup.Group
    var currentOrganization *clerk.Organization
    var profile *model.Profile
    var receiptModel *model.Receipt
    var fileName string
    var pdfFile []byte
    var fileURL string

    // Parallel goroutines for fetching data
    g.Go(func() error {
        var err error
        currentOrganization, err = organization.Get(ctx, input.OrginazationID)
        return err
    })

    g.Go(func() error {
        var err error
        profile, err = r.GetProfileByUserID(userId)
        if err != nil {
            fmt.Println("Could not fetch profile:", err)
        }
        return nil
    })

    if err := g.Wait(); err != nil {
        return false, err
    }

    // Prepare receipt model
    if err := updateProfileImages(profile, currentOrganization); err != nil {
        return false, err
    }

    totalAmount, subtotal, taxAmount := calculateTotalAmount(input, profile.Tax)
    receiptModel = inputToReceiptModel(input, userId, totalAmount, subtotal, taxAmount)
    r.saveReceipt(receiptModel, input)


    // Parallel PDF generation and storage upload
    g.Go(func() error {
        var err error
        fileName, pdfFile, err = r.generatePDF(receiptModel, profile)
        if err != nil {
            return err
        }
        
		if err := r.saveFile(pdfFile, fileName, currentOrganization.ID, userId); err != nil {
			return  err
		}
		fileURL, err = r.getFileURL( currentOrganization.ID, userId, fileName)
		if err != nil {
			return  err
		}
        return err
    })

    if err := g.Wait(); err != nil {
        return false, err
    }

    // Optional: Async WhatsApp message
    if receiptModel.RecipientPhone != "" {
        go func() {
            err := r.sendPDFToWhatsApp(fileURL, fileName, currentOrganization.Name, receiptModel.RecipientPhone)
            if err != nil {
                log.Printf("Failed to send PDF to WhatsApp: %v", err)
            }
        }()
    }

    return true, nil
}

func (r *ReceiptPDFGeneratorResolver) SendReceiptPDFToEmail(ctx context.Context, input model.CreateReceiptPDFGenerator) (bool, error) {
    // Early error checking
    userId, err := contextUtil.UserIdFromContext(ctx)
    if err != nil {
        return false, err
    }

    // Use errgroup for parallel processing
    var currentOrganization *clerk.Organization
    var profile *model.Profile
    var receiptModel *model.Receipt
    var fileName string
    var pdfFile []byte

    currentOrganization, err = organization.Get(ctx, input.OrginazationID)
	if err != nil {
		return false, err
	}
	profile, err = r.GetProfileByUserID(userId)
	if err != nil {
		fmt.Println("Could not fetch profile:", err)
	}

    if err := updateProfileImages(profile, currentOrganization); err != nil {
        return false, err
    }

    totalAmount, subtotal, taxAmount := calculateTotalAmount(input, profile.Tax)
    receiptModel = inputToReceiptModel(input, userId, totalAmount, subtotal, taxAmount)
	r.saveReceipt(receiptModel, input)

    fileName, pdfFile, err = r.generatePDF(receiptModel, profile)
	if err != nil {
		return false, err
	}
    if err := r.saveFile(pdfFile, fileName, currentOrganization.ID, userId); err != nil {
			return false, err
		}
    if receiptModel.RecipientPhone != "" {
			err := emails.SendEmailWithPDF(
				*receiptModel.RecipientEmail,
				"Receipt ",
				"templates/emails/receipt.html",
				map[string]interface{}{
					"OrganizationName": currentOrganization.Name,
					"CustomerName": receiptModel.RecipientName,
				},
				fileName,
				pdfFile,
			)
			if err != nil {
				log.Printf("Failed to send PDF to Email: %v", err)
			}
       
    }

    return true, nil
}
