package receiptPDFGenerator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"text/template"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/cloudFront"
	"github.com/bishal-dd/receipt-generator-backend/helper/search"
)

var (
    templateOnce sync.Once
    cachedTemplate *template.Template
)

func (r *ReceiptPDFGeneratorResolver) saveFile(pdf []byte, fileName string, organizationId string, userId string)( error){
	presignedResp, err := r.httpClient.R().
	Get(fmt.Sprintf("%s/issuePresignedURL?filename=%s&contentType=%s&organizationId=%s&userId=%s",os.Getenv("BACKEND_URL"),fileName, "application/pdf", organizationId, userId ))
	if err != nil {
		return  err
	}
	fmt.Println(presignedResp.String())
	if presignedResp.StatusCode() != http.StatusOK {
		return  err
	}
	// Parse presigned URL response
	type PresignedURLResponse struct {
		URL string `json:"url"`
	}

	var presignedRespData PresignedURLResponse
	if err := json.Unmarshal(presignedResp.Body(), &presignedRespData); err != nil {
		return  err
	}

	presignedURL := presignedRespData.URL
	// Upload the PDF
	uploadResp, err := r.httpClient.R().
		SetHeader("Content-Type", "application/pdf").
		SetBody(pdf).
		Put(presignedURL)	
	if err != nil {
		return   err
	}
	if uploadResp.StatusCode() != http.StatusOK {
		return  fmt.Errorf("failed to upload PDF: %s", uploadResp.String())
	}
	return nil
}	

func (r *ReceiptPDFGeneratorResolver) getFileURL(organizationId, userId, fileName string) (string, error) {
	key := fmt.Sprintf("%s/%s/%s", organizationId, userId, fileName)
	signedURL, err := cloudFront.GetCloudFrontURL(key)
	if err != nil {
		return  "", err
	}
	return  signedURL, nil
}

func (r *ReceiptPDFGeneratorResolver) generatePDF(receipt *model.Receipt, profile *model.Profile) (string, []byte, error) {
	templateOnce.Do(func() {
        templateFile := "templates/receiptTemplate/index.html"
        templateContent, _ := os.ReadFile(templateFile)
        cachedTemplate = template.Must(template.New("receipt").Parse(string(templateContent)))
    })

    var htmlBuffer bytes.Buffer
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
		return  "", nil, fmt.Errorf("gotenberg returned status %d: %s", resp.StatusCode(), resp.String())
	}
	outputFilename := fmt.Sprintf("receipt_%s.pdf", receipt.ID)
	return  outputFilename, resp.Body(), nil
}


func (r *ReceiptPDFGeneratorResolver)  sendPDFToWhatsApp(url string, receiptName string, orginaztionName string, recipientPhone string, receiptId string) error {
	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to": recipientPhone,
		"type": "document",
		"document": map[string]string{
    		"link": url,
			"caption": fmt.Sprintf("Receipt from %s", orginaztionName),
			"filename": receiptName,
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
	receipt := &model.Receipt{
        ID: receiptId,
    }

    
	isReceiptSend := true
	if err := r.db.Model(receipt).Updates(model.UpdateReceipt{IsReceiptSend: &isReceiptSend}).Error; err != nil {
        return  err
    }
    if err := search.UpdateReceiptDocument(r.httpClient, map[string]interface{}{"is_receipt_send": true}, receiptId); err != nil {
        return  err
    }
    return nil
}
