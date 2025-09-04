package receiptPDFGenerator

import (
	"context"
	"errors"
	"fmt"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/contextUtil"
	"github.com/bishal-dd/receipt-generator-backend/helper/emails"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/organization"
	"golang.org/x/sync/errgroup"
)

func (r *ReceiptPDFGeneratorResolver) SendReceiptPDFToWhatsApp(ctx context.Context, input model.SendReceiptPDFToWhatsApp) (bool, error) {
	// Early error checking
	tx := r.db.Begin()
	userId, err := contextUtil.UserIdFromContext(ctx)
	if err != nil {
		return false, err
	}
	user, err := r.GetUserFromDB(ctx, userId)
	if err != nil {
		return false, err
	}
	if err := checkUserMode(user.Mode, user.UseCount); err != nil {
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

	totalAmount, subtotal, taxAmount := calculateTotalAmount(input.Services, profile.Tax)
	receiptModel = whatsAppInputToReceiptModel(input, userId, totalAmount, subtotal, taxAmount)
	encryptedReceiptModel, err := whatsAppInputToEncryptedReceiptModel(input, userId, totalAmount, subtotal, taxAmount, r.publicKeyPEM)
	if err != nil {
		return false, err
	}
	if err := r.saveEncryptedReceipt(encryptedReceiptModel, input.Services, tx); err != nil {
		tx.Rollback()
		return false, err
	}
	if err := r.MinusProductQuantity(input.Services, tx); err != nil {
		tx.Rollback()
		return false, err
	}
	if err := tx.Commit().Error; err != nil {
		return false, err
	}

	// Parallel PDF generation and storage upload
	g.Go(func() error {
		var err error
		fileName, pdfFile, err = r.generatePDF(receiptModel, profile)
		if err != nil {
			return err
		}

		if err := r.saveFile(pdfFile, fileName, currentOrganization.ID, userId); err != nil {
			return err
		}
		fileURL, err = r.getFileURL(currentOrganization.ID, userId, fileName)
		if err != nil {
			return err
		}
		return err
	})

	if err := g.Wait(); err != nil {
		return false, err
	}

	// Optional: Async WhatsApp message
	if receiptModel.RecipientPhone != nil && *receiptModel.RecipientPhone != "" {
		err := r.sendPDFToWhatsApp(fileURL, fileName, currentOrganization.Name, *receiptModel.RecipientPhone, *receiptModel.TotalAmount, encryptedReceiptModel.ID)
		if err != nil {
			return false, err
		}

	} else {
		return false, errors.New("recipient phone is empty")
	}

	return true, nil
}

func (r *ReceiptPDFGeneratorResolver) SendEncryptedReceiptPDFToWhatsAppWithReceiptID(ctx context.Context, receiptId string, orginazationId string, phoneNumber string) (bool, error) {
	// Early error checking
	userId, err := contextUtil.UserIdFromContext(ctx)
	if err != nil {
		return false, err
	}
	user, err := r.GetUserFromDB(ctx, userId)
	if err != nil {
		return false, err
	}
	if err := checkUserMode(user.Mode, user.UseCount); err != nil {
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
		currentOrganization, err = organization.Get(ctx, orginazationId)
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

	receiptModel, err = r.GetEncryptedReceiptFromDB(ctx, receiptId)
	if err != nil {
		return false, err
	}

	// Parallel PDF generation and storage upload
	g.Go(func() error {
		var err error
		fileName, pdfFile, err = r.generatePDF(receiptModel, profile)
		if err != nil {
			return err
		}

		if err := r.saveFile(pdfFile, fileName, currentOrganization.ID, userId); err != nil {
			return err
		}
		fileURL, err = r.getFileURL(currentOrganization.ID, userId, fileName)
		if err != nil {
			return err
		}
		return err
	})

	if err := g.Wait(); err != nil {
		return false, err
	}

	// Optional: Async WhatsApp message
	if phoneNumber != "" {
		err := r.sendPDFToWhatsApp(fileURL, fileName, currentOrganization.Name, phoneNumber, *receiptModel.TotalAmount, receiptModel.ID)
		if err != nil {
			return false, err
		}
	} else {
		return false, errors.New("recipient phone is empty")
	}

	return true, nil
}

func (r *ReceiptPDFGeneratorResolver) SendReceiptPDFToEmail(ctx context.Context, input model.SendReceiptPDFToEmail) (bool, error) {
	// Early error checking
	tx := r.db.Begin()

	userId, err := contextUtil.UserIdFromContext(ctx)
	if err != nil {
		return false, err
	}
	user, err := r.GetUserFromDB(ctx, userId)
	if err != nil {
		return false, err
	}
	if err := checkUserMode(user.Mode, user.UseCount); err != nil {
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
		return false, err
	}

	if err := updateProfileImages(profile, currentOrganization); err != nil {
		return false, err
	}

	totalAmount, subtotal, taxAmount := calculateTotalAmount(input.Services, profile.Tax)
	receiptModel = emailInputToReceiptModel(input, userId, totalAmount, subtotal, taxAmount)
	encryptedReceiptModel, err := emailInputToEncryptedReceiptModel(input, userId, totalAmount, subtotal, taxAmount, r.publicKeyPEM)
	if err != nil {
		return false, err
	}
	if err := r.saveEncryptedReceipt(encryptedReceiptModel, input.Services, tx); err != nil {
		tx.Rollback()
		return false, err
	}
	if err := r.MinusProductQuantity(input.Services, tx); err != nil {
		tx.Rollback()
		return false, err
	}
	if err := tx.Commit().Error; err != nil {
		return false, err
	}
	fileName, pdfFile, err = r.generatePDF(receiptModel, profile)
	if err != nil {
		return false, err
	}
	if err := r.saveFile(pdfFile, fileName, currentOrganization.ID, userId); err != nil {
		return false, err
	}
	if receiptModel.RecipientEmail != nil && *receiptModel.RecipientEmail != "" {
		err := emails.SendEmailWithPDF(
			*receiptModel.RecipientEmail,
			"Receipt",
			"templates/emails/receipt.html",
			map[string]interface{}{
				"OrganizationName": currentOrganization.Name,
				"CustomerName":     receiptModel.RecipientName,
			},
			fileName,
			pdfFile,
		)
		if err != nil {
			return false, err
		}
		receipt := &model.EncryptedReceipt{
			ID: encryptedReceiptModel.ID,
		}

		isReceiptSend := true
		if err := r.db.Model(receipt).Updates(model.UpdateEncryptedReceipt{IsReceiptSend: &isReceiptSend}).Error; err != nil {
			return false, err
		}
		// if err := search.UpdateReceiptDocument(r.httpClient, map[string]interface{}{"is_receipt_send": true}, receiptModel.ID); err != nil {
		//     return  false, err
		// }
	} else {
		return false, errors.New("recipient email is empty")
	}

	return true, nil
}

func (r *ReceiptPDFGeneratorResolver) SendEncryptedReceiptPDFToEmailWithReceiptID(ctx context.Context, receiptId string, orginazationId string, email string) (bool, error) {
	// Early error checking
	userId, err := contextUtil.UserIdFromContext(ctx)
	if err != nil {
		return false, err
	}
	user, err := r.GetUserFromDB(ctx, userId)
	if err != nil {
		return false, err
	}
	if err := checkUserMode(user.Mode, user.UseCount); err != nil {
		return false, err
	}
	// Use errgroup for parallel processing
	var currentOrganization *clerk.Organization
	var profile *model.Profile
	var receiptModel *model.Receipt
	var fileName string
	var pdfFile []byte

	currentOrganization, err = organization.Get(ctx, orginazationId)
	if err != nil {
		return false, err
	}
	profile, err = r.GetProfileByUserID(userId)
	if err != nil {
		return false, err
	}

	if err := updateProfileImages(profile, currentOrganization); err != nil {
		return false, err
	}

	receiptModel, err = r.GetEncryptedReceiptFromDB(ctx, receiptId)
	if err != nil {
		return false, err
	}
	fileName, pdfFile, err = r.generatePDF(receiptModel, profile)
	if err != nil {
		return false, err
	}
	if err := r.saveFile(pdfFile, fileName, currentOrganization.ID, userId); err != nil {
		return false, err
	}
	err = emails.SendEmailWithPDF(
		email,
		"Receipt",
		"templates/emails/receipt.html",
		map[string]interface{}{
			"OrganizationName": currentOrganization.Name,
			"CustomerName":     receiptModel.RecipientName,
		},
		fileName,
		pdfFile,
	)
	if err != nil {
		return false, err
	}

	receipt := &model.EncryptedReceipt{
		ID: receiptModel.ID,
	}

	isReceiptSend := true
	if err := r.db.Model(receipt).Updates(model.UpdateEncryptedReceipt{IsReceiptSend: &isReceiptSend}).Error; err != nil {
		return false, err
	}
	// if err := search.UpdateReceiptDocument(r.httpClient, map[string]interface{}{"is_receipt_send": true}, receiptModel.ID); err != nil {
	//     return  false, err
	// }

	return true, nil
}

func (r *ReceiptPDFGeneratorResolver) DownloadReceiptPDF(ctx context.Context, input model.DownloadPDF) (string, error) {
	// Early error checking
	tx := r.db.Begin()
	userId, err := contextUtil.UserIdFromContext(ctx)
	if err != nil {
		return "", err
	}
	user, err := r.GetUserFromDB(ctx, userId)
	if err != nil {
		return "", err
	}
	if err := checkUserMode(user.Mode, user.UseCount); err != nil {
		return "", err
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
		return "", err
	}

	// Prepare receipt model
	if err := updateProfileImages(profile, currentOrganization); err != nil {
		return "", err
	}

	totalAmount, subtotal, taxAmount := calculateTotalAmount(input.Services, profile.Tax)
	receiptModel = downloadInputToReceiptModel(input, userId, totalAmount, subtotal, taxAmount)
	encryptedReceiptModel, err := downlaodInputToEncryptedReceiptModel(input, userId, totalAmount, subtotal, taxAmount, r.publicKeyPEM)
	if err != nil {
		return "", err
	}
	if err := r.saveEncryptedReceipt(encryptedReceiptModel, input.Services, tx); err != nil {
		tx.Rollback()
		return "", err
	}
	if err := r.MinusProductQuantity(input.Services, tx); err != nil {
		tx.Rollback()
		return "", err
	}
	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	// Parallel PDF generation and storage upload
	g.Go(func() error {
		var err error
		fileName, pdfFile, err = r.generatePDF(receiptModel, profile)
		if err != nil {
			return err
		}

		if err := r.saveFile(pdfFile, fileName, currentOrganization.ID, userId); err != nil {
			return err
		}
		fileURL, err = r.getFileURL(currentOrganization.ID, userId, fileName)
		if err != nil {
			return err
		}
		return err
	})

	if err := g.Wait(); err != nil {
		return "", err
	}

	return fileURL, nil
}

func (r *ReceiptPDFGeneratorResolver) DownloadEncryptedReceiptPDFWithReceiptID(ctx context.Context, receiptId string, orginazationId string) (string, error) {
	// Early error checking
	userId, err := contextUtil.UserIdFromContext(ctx)
	if err != nil {
		return "", err
	}
	user, err := r.GetUserFromDB(ctx, userId)
	if err != nil {
		return "", err
	}
	if err := checkUserMode(user.Mode, user.UseCount); err != nil {
		return "", err
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
		currentOrganization, err = organization.Get(ctx, orginazationId)
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
		return "", err
	}

	// Prepare receipt model
	if err := updateProfileImages(profile, currentOrganization); err != nil {
		return "", err
	}

	receiptModel, err = r.GetEncryptedReceiptFromDB(ctx, receiptId)
	if err != nil {
		return "", err
	}
	// Parallel PDF generation and storage upload
	g.Go(func() error {
		var err error
		fileName, pdfFile, err = r.generatePDF(receiptModel, profile)
		if err != nil {
			return err
		}

		if err := r.saveFile(pdfFile, fileName, currentOrganization.ID, userId); err != nil {
			return err
		}
		fileURL, err = r.getFileURL(currentOrganization.ID, userId, fileName)
		if err != nil {
			return err
		}
		return err
	})

	if err := g.Wait(); err != nil {
		return "", err
	}

	return fileURL, nil
}

func (r *ReceiptPDFGeneratorResolver) SaveReceipt(ctx context.Context, input model.DownloadPDF) (bool, error) {
	tx := r.db.Begin()
	userId, err := contextUtil.UserIdFromContext(ctx)
	if err != nil {
		return false, err
	}
	var profile *model.Profile
	var receiptModel *model.Receipt

	profile, err = r.GetProfileByUserID(userId)
	if err != nil {
		fmt.Println("Could not fetch profile:", err)
	}

	totalAmount, subtotal, taxAmount := calculateTotalAmount(input.Services, profile.Tax)
	receiptModel = downloadInputToReceiptModel(input, userId, totalAmount, subtotal, taxAmount)
	if err := r.saveReceipt(receiptModel, input.Services, tx); err != nil {
		tx.Rollback()
		return false, err
	}
	if err := r.MinusProductQuantity(input.Services, tx); err != nil {
		tx.Rollback()
		return false, err
	}
	if err := tx.Commit().Error; err != nil {
		return false, err
	}

	return true, nil
}
