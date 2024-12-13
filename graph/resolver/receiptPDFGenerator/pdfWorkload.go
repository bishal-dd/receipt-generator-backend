package receiptPDFGenerator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/cloudFront"
)


func (r *ReceiptPDFGeneratorResolver) getFileURL(pdf []byte, fileName string, organizationId string, userId string)( string, error){
	presignedResp, err := r.httpClient.R().
	Get(fmt.Sprintf("http://localhost:8080/issuePresignedURL?filename=%s&contentType=%s&organizationId=%s&userId=%s",fileName, "application/pdf", organizationId, userId ))
	if err != nil {
		return  "", err
	}

	if presignedResp.StatusCode() != http.StatusOK {
		return  "", err
	}
	// Parse presigned URL response
	type PresignedURLResponse struct {
		URL string `json:"url"`
	}

	var presignedRespData PresignedURLResponse
	if err := json.Unmarshal(presignedResp.Body(), &presignedRespData); err != nil {
		return  "", err
	}

	presignedURL := presignedRespData.URL
	// Upload the PDF
	uploadResp, err := r.httpClient.R().
		SetHeader("Content-Type", "application/pdf").
		SetBody(pdf).
		Put(presignedURL)	
	if err != nil {
		return   "", err
	}
	if uploadResp.StatusCode() != http.StatusOK {
		return  "", fmt.Errorf("failed to upload PDF: %s", uploadResp.String())
	}
	key := fmt.Sprintf("%s/%s/%s", organizationId, userId, fileName)
	signedURL, err := cloudFront.GetCloudFrontURL(key)
	if err != nil {
		return  "", err
	}
	return  signedURL, nil
}

func (r *ReceiptPDFGeneratorResolver) generatePDF(receipt *model.Receipt, profile *model.Profile) (string, []byte, error) {
	templateFile := "templates/receiptTemplate/index.html"
	templateContent, err := os.ReadFile(templateFile)
	if err != nil {
		return  "", nil, err
	}
	tmpl, err := template.New("receipt").Parse(string(templateContent))
	if err != nil {
		return  "", nil, err
	}
	var htmlBuffer bytes.Buffer
	data := struct {
		Receipt *model.Receipt
        Profile *model.Profile
	}{
		Receipt: receipt,
        Profile: profile,
	}
	if err := tmpl.Execute(&htmlBuffer, data); err != nil {
		return "", nil, err
	}
	resp, err := r.httpClient.R().
		SetHeader("Content-Type", "multipart/form-data").
		SetFileReader("files", "index.html", bytes.NewReader(htmlBuffer.Bytes())). 
		SetFormData(map[string]string{
			"index": "index.html", 
		}).
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


func (r *ReceiptPDFGeneratorResolver)  sendPDFToWhatsApp(url string, receiptName string, orginaztionName string,) error {
	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to": "97517959259",
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

    return nil
}
