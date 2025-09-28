package receiptFiles

import (
	"bytes"
	"context"
	"io"

	"regexp"

	"github.com/99designs/gqlgen/graphql"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

func (r *ReceiptFilesResolver) VerifyReceiptFile(ctx context.Context, file graphql.Upload) (bool, error) {
	// 1. Read file bytes
	var buf bytes.Buffer
	_, err := io.Copy(&buf, file.File)
	if err != nil {
		return false, err
	}
	pdfBytes := buf.Bytes()

	// 2. Load PDF
	pdfReader, err := model.NewPdfReader(bytes.NewReader(pdfBytes))
	if err != nil {
		return false, err
	}

	// 3. Extract receipt ID
	receiptID := ""
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return false, err
	}

	for i := 1; i <= numPages; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			continue
		}

		ex, err := extractor.New(page)
		if err != nil {
			continue
		}

		text, err := ex.ExtractText()
		if err != nil {
			continue
		}

		id := extractReceiptIDFromText(text)
		if id != "" {
			receiptID = id
			break
		}
	}

	if receiptID == "" {
		return false, nil // receipt ID not found
	}

	// 4. Check in database
	var exists bool
	err = r.db.
		Raw(`SELECT EXISTS(SELECT 1 FROM receipt_files WHERE id = ?)`, receiptID).
		Scan(&exists).Error
	if err != nil {
		return false, err
	}

	return exists, nil
}

// extractReceiptIDFromText uses regex to find the ID
func extractReceiptIDFromText(text string) string {
	re := regexp.MustCompile(`RECEIPT_ID:([a-f0-9-]+)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
