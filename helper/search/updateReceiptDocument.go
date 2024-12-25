package search

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
)

func UpdateReceiptDocument (httpClient *resty.Client, input map[string]interface{}, id string ) error {
	typeSenseURL := os.Getenv("TYPESENSE_URL")
	if typeSenseURL == "" {
		return fmt.Errorf("TYPESENSE_URL is not set")
	}
	resp, err := httpClient.R().
	SetHeader("Content-Type", "application/json").
	SetHeader("X-TYPESENSE-API-KEY", os.Getenv("TYPESENSE_API_KEY")).
	SetBody(input).
	Patch(fmt.Sprintf("%s/collections/receipts/documents/%s", typeSenseURL, id ) )
	if err != nil {
		return  err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("typesense returned status %d: %s", resp.StatusCode(), resp.String())
	}
	return  nil
}