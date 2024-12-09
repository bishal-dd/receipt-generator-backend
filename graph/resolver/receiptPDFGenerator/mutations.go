package receiptPDFGenerator

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/contextUtil"
	"github.com/bishal-dd/receipt-generator-backend/helper/ids"
	"github.com/jung-kurt/gofpdf"
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

	if err := generatePDF(receiptModel); err != nil {
        return false, err
    }
    return true, nil
}

func generatePDF(receipt *model.Receipt) error {
    // Read HTML template from file
    templateFile := "templates/receipt_template.html"
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
    }{
        Receipt: receipt,
    }
    if err := tmpl.Execute(&htmlBuffer, data); err != nil {
        return fmt.Errorf("error rendering HTML template: %w", err)
    }

    // Create a new PDF
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    pdf.SetFont("Arial", "", 12)

    // Manually parse and write content
    pdf.SetLeftMargin(10)
    pdf.SetRightMargin(10)

    // Title
    pdf.SetFont("Arial", "B", 16)
    pdf.CellFormat(190, 10, fmt.Sprintf("Receipt: %s", receipt.ReceiptName), "", 1, "C", false, 0, "")
    pdf.Ln(10)

    // Receipt Details
    pdf.SetFont("Arial", "", 12)
    pdf.CellFormat(190, 7, fmt.Sprintf("Date: %s", receipt.Date), "", 1, "L", false, 0, "")
    pdf.CellFormat(190, 7, fmt.Sprintf("Recipient: %s (%d)", receipt.RecipientName, receipt.RecipientPhone), "", 1, "L", false, 0, "")
    pdf.CellFormat(190, 7, fmt.Sprintf("Total Amount: %v", receipt.TotalAmount), "", 1, "L", false, 0, "")
    
    pdf.Ln(10)

    // Table Header
    pdf.SetFont("Arial", "B", 12)
    pdf.CellFormat(47.5, 7, "Description", "1", 0, "C", false, 0, "")
    pdf.CellFormat(47.5, 7, "Rate", "1", 0, "C", false, 0, "")
    pdf.CellFormat(47.5, 7, "Quantity", "1", 0, "C", false, 0, "")
    pdf.CellFormat(47.5, 7, "Amount", "1", 1, "C", false, 0, "")

    // Table Rows
    pdf.SetFont("Arial", "", 12)
    for _, service := range receipt.Services {
        pdf.CellFormat(47.5, 7, service.Description, "1", 0, "L", false, 0, "")
        pdf.CellFormat(47.5, 7, fmt.Sprintf("%v", service.Rate), "1", 0, "R", false, 0, "")
        pdf.CellFormat(47.5, 7, fmt.Sprintf("%v", service.Quantity), "1", 0, "R", false, 0, "")
        pdf.CellFormat(47.5, 7, fmt.Sprintf("%v", service.Amount), "1", 1, "R", false, 0, "")
    }

    // Output the PDF
    outputFilename := fmt.Sprintf("receipt_%s.pdf", receipt.ID)
    if err := pdf.OutputFileAndClose(outputFilename); err != nil {
        return fmt.Errorf("error generating PDF: %w", err)
    }

    fmt.Println("PDF generated successfully:", outputFilename)
    return nil
}