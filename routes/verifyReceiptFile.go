package routes

import (
	"bytes"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	pdf "github.com/ledongthuc/pdf"
	"gorm.io/gorm"
)

// VerifyReceiptFileHandler validates an uploaded receipt PDF.
func VerifyReceiptFileHandler(c *gin.Context, db *gorm.DB) {
	respond := func(valid bool, message string, receiptID ...string) {
		resp := gin.H{"valid": valid, "message": message}
		if len(receiptID) > 0 {
			resp["id"] = receiptID[0]
		}
		c.JSON(http.StatusOK, resp)
	}

	// Parse multipart form file
	file, err := c.FormFile("file")
	if err != nil {
		respond(false, "file is required")
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		respond(false, "cannot open file")
		return
	}
	defer openedFile.Close()

	// Read file bytes
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(openedFile); err != nil {
		respond(false, "cannot read file")
		return
	}
	pdfBytes := buf.Bytes()

	// Parse PDF
	reader, err := pdf.NewReader(bytes.NewReader(pdfBytes), int64(len(pdfBytes)))
	if err != nil {
		respond(false, "invalid PDF")
		return
	}

	// Extract RECEIPT_ID from all pages
	var receiptID string
	for i := 1; i <= reader.NumPage(); i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}
		content, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}

		if id := extractReceiptID(content); id != "" {
			receiptID = id
			break
		}
	}

	if receiptID == "" {
		respond(false, "Receipt ID not found")
		return
	}

	// Check in database
	var exists bool
	if err := db.Raw(`SELECT EXISTS(SELECT 1 FROM receipt_files WHERE id = ?)`, receiptID).
		Scan(&exists).Error; err != nil {
		respond(false, "Database query failed")
		return
	}
	if !exists {
		respond(false, "Receipt ID not found")
		return
	}

	// Success
	respond(true, "Receipt is valid", receiptID)
}

// extractReceiptID tries to find the hidden RECEIPT_ID string in text.
func extractReceiptID(text string) string {
	re := regexp.MustCompile(`RECEIPT_ID:\s*([a-f0-9-]+)`)
	if matches := re.FindStringSubmatch(text); len(matches) > 1 {
		return matches[1]
	}
	return ""
}
