package search

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/go-resty/resty/v2"
)

func SearchReceiptDocuments(httpClient *resty.Client, userId string, page int, year *int, date *string, dateRange []string) (*model.SearchReceipt, error) {
	typeSenseURL := os.Getenv("TYPESENSE_URL")
	if typeSenseURL == "" {
		return nil, fmt.Errorf("TYPESENSE_URL is not set")
	}

	filters := "user_id:=" + userId // Initialize the filters with user_id condition.

	if year != nil {
		filters += fmt.Sprintf(" && year:=%d", *year) // Add `year` condition.
	}
	if date != nil {
		filters += fmt.Sprintf(" && date:=%s", *date) // Add `date` condition.
	}
	if dateRange != nil {
		parsedDate1, err := time.Parse("2006-01-02", dateRange[0])
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %v", err)
		}
		parsedDate2, err := time.Parse("2006-01-02", dateRange[1])
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %v", err)
		}
		parsedDate1 = parsedDate1.AddDate(0, 0, -1)
		unixDate1 := parsedDate1.Unix()
		unixDate2 := parsedDate2.Unix()
		
		filters += fmt.Sprintf(" && date_unix:[%d..%d]", unixDate1, unixDate2) // Add `dateRange` condition.
	}

	resp, err := httpClient.R().
		SetHeader("X-TYPESENSE-API-KEY", os.Getenv("TYPESENSE_API_KEY")).
		Get(fmt.Sprintf("%s/collections/receipts/documents/search?q=*&page=%d&per_page=10&filter_by=%s", typeSenseURL, page, filters))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("typesense returned status %d: %s", resp.StatusCode(), resp.String())
	}

	response, err := ParseSearchResponse(resp.Body())
	if err != nil {
		return nil, err
	}

	receipts := ExtractReceipts(response)
	return &model.SearchReceipt{
		Receipts:   receipts,
		TotalCount: response.OutOf,
		FoundCount: response.Found,
	}, nil
}
