package receiptPDFGenerator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"text/template"
	"time"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/cloudFront"
	"github.com/bishal-dd/receipt-generator-backend/helper/stringUtil"
	"gorm.io/gorm"
)

var (
	templateOnce   sync.Once
	cachedTemplate *template.Template
)

func (r *ReceiptPDFGeneratorResolver) saveFile(pdf []byte, fileName string, organizationId string, userId string) error {
	presignedResp, err := r.httpClient.R().
		Get(fmt.Sprintf("%s/issuePresignedURL?filename=%s&contentType=%s&organizationId=%s&userId=%s", os.Getenv("BACKEND_URL"), fileName, "application/pdf", organizationId, userId))
	if err != nil {
		return err
	}
	if presignedResp.StatusCode() != http.StatusOK {
		return err
	}
	// Parse presigned URL response
	type PresignedURLResponse struct {
		URL string `json:"url"`
	}

	var presignedRespData PresignedURLResponse
	if err := json.Unmarshal(presignedResp.Body(), &presignedRespData); err != nil {
		return err
	}

	presignedURL := presignedRespData.URL
	// Upload the PDF
	uploadResp, err := r.httpClient.R().
		SetHeader("Content-Type", "application/pdf").
		SetBody(pdf).
		Put(presignedURL)
	if err != nil {
		return err
	}
	if uploadResp.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to upload PDF: %s", uploadResp.String())
	}
	user := &model.User{}
	if err := r.db.Model(user).Where("id = ?", userId).Update("use_count", gorm.Expr("use_count + 1")).Error; err != nil {
		return fmt.Errorf("failed to increment use_count for user %s: %w", userId, err)
	}
	return nil
}

func (r *ReceiptPDFGeneratorResolver) getFileURL(organizationId, userId, fileName string) (string, error) {
	key := fmt.Sprintf("%s/%s/%s", organizationId, userId, fileName)
	signedURL, err := cloudFront.GetCloudFrontURL(key)
	if err != nil {
		return "", err
	}
	return signedURL, nil
}

func (r *ReceiptPDFGeneratorResolver) generatePDF(receipt *model.Receipt, profile *model.Profile) (string, []byte, error) {
	templateOnce.Do(func() {
		templateFile := "templates/receiptTemplate/index.html"
		templateContent, _ := os.ReadFile(templateFile)
		cachedTemplate = template.Must(template.New("receipt").Parse(string(templateContent)))
	})

	parsedDate, err := time.Parse(time.RFC3339, receipt.Date)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse receipt date: %w", err)
	}
	receipt.Date = parsedDate.Format("January 02 2006")
	var htmlBuffer bytes.Buffer
	for i, svc := range receipt.Services {
		fmt.Printf("Service %d: %+v\n", i, *svc)
	}

	data := struct {
		Receipt *model.Receipt
		Profile *model.Profile
	}{
		Receipt: receipt,
		Profile: profile,
	}

	if err := cachedTemplate.Execute(&htmlBuffer, data); err != nil {
		return "", nil, err
	}

	// Consider using a Gotenberg connection pool
	resp, err := r.httpClient.R().
		SetHeader("Content-Type", "multipart/form-data").
		SetFileReader("files", "index.html", bytes.NewReader(htmlBuffer.Bytes())).
		SetFormData(map[string]string{"index": "index.html"}).
		Post(os.Getenv("GOTENBERY_URL"))

	if err != nil {
		return "", nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return "", nil, fmt.Errorf("gotenberg returned status %d: %s", resp.StatusCode(), resp.String())
	}
	outputFilename := fmt.Sprintf("receipt_of_%v_from_%s_%s.pdf", *receipt.TotalAmount, stringUtil.ReplaceSpaceWithUnderscore(*profile.CompanyName), receipt.ID)
	return outputFilename, resp.Body(), nil
}

func (r *ReceiptPDFGeneratorResolver) sendPDFToWhatsApp(url string, receiptName string, orginaztionName string, recipientPhone string, receiptAmount float64, receiptId string) error {
	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                recipientPhone,
		"type":              "template",
		"template": map[string]interface{}{
			"name": "receipt_template", // Replace with your approved template name
			"language": map[string]string{
				"code": "en_US", // Adjust as necessary
			},
			"components": []map[string]interface{}{
				{
					"type": "header",
					"parameters": []map[string]interface{}{
						{
							"type": "document",
							"document": map[string]interface{}{
								"link":     url, // URL to the document (PDF)
								"filename": receiptName,
							},
						},
					},
				},
				{
					"type": "body",
					"parameters": []map[string]interface{}{
						{"type": "text", "text": receiptAmount},
						{"type": "text", "text": orginaztionName},
						{"type": "text", "text": "receipt"},
					},
				},
			},
		},
	}

	resp, err := r.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(os.Getenv("WHATSAPP_ACCESS_TOKEN")).
		SetBody(payload).
		Post(fmt.Sprintf("https://graph.facebook.com/v21.0/%s/messages", os.Getenv("WHATSAPP_PHONE_NUMBER_ID")))

	if err != nil {
		return fmt.Errorf("error sending WhatsApp message: %w", err)
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("WhatsApp API request failed with status %d: %s",
			resp.StatusCode(),
			string(resp.Body()))
	}
	fmt.Printf("%s", receiptId)
	encryptedReceipt := &model.EncryptedReceipt{
		ID: receiptId,
	}
	isReceiptSend := true
	if err := r.db.Model(encryptedReceipt).Updates(model.UpdateEncryptedReceipt{IsReceiptSend: &isReceiptSend}).Error; err != nil {
		return err
	}
	// if err := search.UpdateReceiptDocument(r.httpClient, map[string]interface{}{"is_receipt_send": true}, receiptId); err != nil {
	//     return  err
	// }

	return nil
}
