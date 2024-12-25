package search

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/go-resty/resty/v2"
)

func SearchReceiptDocuments (httpClient *resty.Client, userId string, page int )(*model.SearchReceipt, error) {
	typeSenseURL := os.Getenv("TYPESENSE_URL")
	if typeSenseURL == "" {
		return nil, fmt.Errorf("TYPESENSE_URL is not set")
	}
	resp, err := httpClient.R().
	SetHeader("X-TYPESENSE-API-KEY", os.Getenv("TYPESENSE_API_KEY")).
	Get(fmt.Sprintf("%s/collections/receipts/documents/search?q=%s&query_by=user_id&page=%d&per_page=10&filter_by=date:2024", typeSenseURL, userId, page ) )
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
	fmt.Println(receipts)
	return &model.SearchReceipt{
		Receipts: receipts,
		TotalCount: response.OutOf,
		FoundCount: response.Found,
	}, nil
}