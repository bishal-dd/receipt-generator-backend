package search

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/go-resty/resty/v2"
)

func SearchReceiptDocuments (httpClient *resty.Client, userId string, page int, year *int, date *string, dateRange []string )(*model.SearchReceipt, error) {
	typeSenseURL := os.Getenv("TYPESENSE_URL")
	if typeSenseURL == "" {
		return nil, fmt.Errorf("TYPESENSE_URL is not set")
	}
	filters := ""

	if year != nil {
		filters = fmt.Sprintf("&filter_by=year:=%d", *year)
	}
	if date != nil {
		filters = fmt.Sprintf("&filter_by=date:=%s", *date)
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
		unixDate1 := parsedDate1.Unix()
		unixDate2 := parsedDate2.Unix()
		filters = fmt.Sprintf("&filter_by=date_unix:[%d..%d]", unixDate1, unixDate2)
	}
	resp, err := httpClient.R().
	SetHeader("X-TYPESENSE-API-KEY", os.Getenv("TYPESENSE_API_KEY")).
	Get(fmt.Sprintf("%s/collections/receipts/documents/search?q=%s&query_by=user_id&page=%d&per_page=10%s", typeSenseURL, userId, page, filters) )
	if err != nil {
		return nil,err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil,fmt.Errorf("typesense returned status %d: %s", resp.StatusCode(), resp.String())
	}
	response, err := ParseSearchResponse(resp.Body())
	if err != nil {
		return nil, err
	}
	receipts := ExtractReceipts(response)
	return &model.SearchReceipt{
		Receipts: receipts,
		TotalCount: response.OutOf,
		FoundCount: response.Found,
	}, nil
}