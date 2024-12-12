package receiptPDFGenerator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/cloudFront"
	"github.com/bishal-dd/receipt-generator-backend/helper/contextUtil"
	"github.com/bishal-dd/receipt-generator-backend/helper/ids"
	"github.com/clerk/clerk-sdk-go/v2/organization"
)


func (r *ReceiptPDFGeneratorResolver) CreateReceiptPDFGenerator(ctx context.Context, input model.CreateReceiptPDFGenerator) (bool, error) {
    // Get user ID from context
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
        // You might want to handle this differently
        // For now, we'll just log it and continue
        fmt.Println("Could not fetch profile:", err)
    }

	if profile.SignatureImage != nil && *profile.SignatureImage != "" {
		signedURL, err := cloudFront.GetCloudFrontURL(*profile.SignatureImage)
		if err != nil {
			return false, err
		}
		profile.SignatureImage = &signedURL
	}
	profile.CompanyName = organization.Name
	if (organization.HasImage){
		profile.LogoImage = organization.ImageURL
	}
	

	subtotal := 0.0
    for _, serviceInput := range input.Services {
		subtotal += serviceInput.Amount
    }

	taxRate := profile.Tax / 100
	taxAmount := subtotal * taxRate
	totalAmount := subtotal + taxAmount

    // Create Receipt
    receiptInput := input
    receiptModel := &model.Receipt{
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
	parsedDate, err := time.Parse(time.RFC3339, receiptModel.Date)
	if err != nil {
		tx.Rollback()
		return false, fmt.Errorf("invalid date format in receipt: %w", err)
	}
	receiptModel.Date = parsedDate.Format("2 January 2006")
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
   
	if err := r.generatePDF(receiptModel, profile, userId, organization.ID); err != nil {
        return false, err
    }
    return true, nil
}

func (r *ReceiptPDFGeneratorResolver) generatePDF(receipt *model.Receipt, profile *model.Profile, userId string, organizationId string) error {
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

	gotenbergURL := "https://gotenberg-production-70d3.up.railway.app/forms/chromium/convert/html"

	resp, err := r.httpClient.R().
		SetHeader("Content-Type", "multipart/form-data").
		SetFileReader("files", "index.html", bytes.NewReader(htmlBuffer.Bytes())). 
		SetFormData(map[string]string{
			"index": "index.html", 
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
	// if err := os.WriteFile(outputFilename, resp.Body(), 0644); err != nil {
	// 	return fmt.Errorf("error saving PDF to file: %w", err)
	// }

	


	// Fetch the presigned URL
	presignedResp, err := r.httpClient.R().
		Get(fmt.Sprintf("http://localhost:8080/issuePresignedURL?filename=%s&contentType=%s&organizationId=%s&userId=%s",outputFilename, "application/pdf", organizationId, userId ))
	if err != nil {
		return fmt.Errorf("error getting presigned URL: %w", err)
	}

	if presignedResp.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to get presigned URL: %s", presignedResp.String())
	}
	// Parse presigned URL response
	type PresignedURLResponse struct {
		URL string `json:"url"`
	}

	var presignedRespData PresignedURLResponse
	if err := json.Unmarshal(presignedResp.Body(), &presignedRespData); err != nil {
		return fmt.Errorf("error parsing presigned URL: %w", err)
	}

	presignedURL := presignedRespData.URL
	// Upload the PDF
	uploadResp, err := r.httpClient.R().
		SetHeader("Content-Type", "application/pdf").
		SetBody(resp.Body()).
		Put(presignedURL)	
	if err != nil {
		return fmt.Errorf("error uploading PDF: %w", err)
	}
	if uploadResp.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to upload PDF: %s", uploadResp.String())
	}
	key := fmt.Sprintf("%s/%s/%s", organizationId, userId, outputFilename)
	signedURL, err := cloudFront.GetCloudFrontURL(key)
	if err != nil {
		return fmt.Errorf("failed to get cloudfront url: %v", err)
	}
	if receipt.RecipientPhone != 0 {
        err = r.sendPDFToWhatsApp(signedURL,outputFilename, receipt)
        if err != nil {
            log.Printf("Failed to send PDF to WhatsApp: %v", err)
			return fmt.Errorf("failed to send pdf to whatsapp: %v", err)

        }
    }

	return nil
}


func (r *ReceiptPDFGeneratorResolver)  sendPDFToWhatsApp(url string, receiptName string, receipt *model.Receipt,) error {
    // Encode PDF to base64
    // base64PDF := base64.StdEncoding.EncodeToString(pdfBytes)
    // Prepare the API request payload
	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to": "97517959259",
		"type": "document",
		"document": map[string]string{
    		"link": url,
			"caption": "Receipt",
			"filename": receiptName,
		},
    }

    // Create Resty client
   

    // Send the request
    resp, err := r.httpClient.R().
        SetHeader("Content-Type", "application/json").
        SetAuthToken(os.Getenv("WHATSAPP_ACCESS_TOKEN")).
        SetBody(payload).
        Post(fmt.Sprintf("https://graph.facebook.com/v21.0/%s/messages", os.Getenv("WHATSAPP_PHONE_NUMBER_ID")))

    if err != nil {
        return fmt.Errorf("error sending WhatsApp message: %w", err)
    }
    fmt.Printf("Response Body: %s\n", string(resp.Body()))    // Check response
    if resp.StatusCode() != 200 {
        return fmt.Errorf("WhatsApp API request failed with status %d: %s", 
            resp.StatusCode(), 
            string(resp.Body()))
    }

    return nil
}
