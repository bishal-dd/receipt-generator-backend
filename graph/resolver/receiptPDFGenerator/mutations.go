package receiptPDFGenerator

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/contextUtil"
	"github.com/bishal-dd/receipt-generator-backend/helper/ids"
	"github.com/go-resty/resty/v2"
)


func (r *ReceiptPDFGeneratorResolver) CreateReceiptPDFGenerator(ctx context.Context, input model.CreateReceiptPDFGenerator) (bool, error) {
    // Get user ID from context
    userId, err := contextUtil.UserIdFromContext(ctx)
    if err != nil {
        return false, err
    }

    // Create Receipt
    receiptInput := input
    receiptModel := &model.Receipt{
        ID:        ids.UUID(),
        UserID:    userId,
        CreatedAt: time.Now().Format(time.RFC3339),
		ReceiptName: receiptInput.ReceiptName,
		RecipientPhone: receiptInput.RecipientPhone,
		RecipientName: receiptInput.RecipientName,
		Amount: receiptInput.Amount,
		TransactionNo: receiptInput.TransactionNo,
		Date: receiptInput.Date,
		TotalAmount: receiptInput.TotalAmount,
		Services:  make([]*model.Service, 0),
    }

    // Begin a transaction
    tx := r.db.Begin()
    if tx.Error != nil {
        return false, tx.Error
    }

    // Create Receipt
    if err := tx.Create(receiptModel).Error; err != nil {
        tx.Rollback()
        return false, err
    }

    // Create Services
    for _, serviceInput := range input.Services {
        serviceModel := &model.Service{
            ID:         ids.UUID(),
            ReceiptID:  receiptModel.ID,
            CreatedAt:  time.Now().Format(time.RFC3339),
            Description: serviceInput.Description,
            Rate:      serviceInput.Rate,
            Quantity:   serviceInput.Quantity,
			Amount:    serviceInput.Amount,
        }

        if err := tx.Create(serviceModel).Error; err != nil {
            tx.Rollback()
            return false, err
        }
		receiptModel.Services = append(receiptModel.Services, serviceModel)
    }

    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        return false, err
    }

    // Get profile (optional)
    profile, err := r.GetProfileByUserID(userId)
    if err != nil {
        // You might want to handle this differently
        // For now, we'll just log it and continue
        fmt.Println("Could not fetch profile:", err)
    }
	fmt.Print("Profile: ", profile)

	if err := generatePDF(receiptModel, profile); err != nil {
        return false, err
    }
    return true, nil
}

func generatePDF(receipt *model.Receipt, profile *model.Profile) error {
	// Read HTML template from file
	templateFile := "templates/receiptTemplate/index.html"
	templateContent, err := os.ReadFile(templateFile)
	if err != nil {
		return fmt.Errorf("error reading HTML template: %w", err)
	}

	// Parse the template
	tmpl, err := template.New("receipt").Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("error parsing HTML template: %w", err)
	}

	// Render the template into a buffer
	var htmlBuffer bytes.Buffer
	data := struct {
		Receipt *model.Receipt
        Profile *model.Profile
	}{
		Receipt: receipt,
        Profile: profile,
	}
	if err := tmpl.Execute(&htmlBuffer, data); err != nil {
		return fmt.Errorf("error rendering HTML template: %w", err)
	}

	// Use Resty to send the HTTP request to Gotenberg
	client := resty.New()
	gotenbergURL := "https://gotenberg-production-70d3.up.railway.app/forms/chromium/convert/html"

	resp, err := client.R().
		SetHeader("Content-Type", "multipart/form-data").
		SetFileReader("files", "index.html", bytes.NewReader(htmlBuffer.Bytes())). // Add the rendered HTML as a file
		SetFormData(map[string]string{
			"index": "index.html", // Specify the index file name
		}).
		Post(gotenbergURL)

	if err != nil {
		return fmt.Errorf("error sending request to Gotenberg: %w", err)
	}

	// Check response status
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("gotenberg returned status %d: %s", resp.StatusCode(), resp.String())
	}

	// Save the PDF to a file
	outputFilename := fmt.Sprintf("receipt_%s.pdf", receipt.ID)
	if err := os.WriteFile(outputFilename, resp.Body(), 0644); err != nil {
		return fmt.Errorf("error saving PDF to file: %w", err)
	}

	fmt.Println("PDF generated successfully:", outputFilename)
	return nil
}
