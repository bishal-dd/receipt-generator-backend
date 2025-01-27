package product

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/go-resty/resty/v2"
)

// SearchResponse represents the top-level response from TypeSense
type SearchResponse struct {
    FacetCounts []interface{} `json:"facet_counts"`
    Found       int           `json:"found"`
    Hits        []Hit         `json:"hits"`
    OutOf       int          `json:"out_of"`
    Page        int          `json:"page"`
    SearchTimeMs int         `json:"search_time_ms"`
}

// Hit represents each document hit in the search results
type Hit struct {
    Document     model.Product    `json:"document"`
    Highlight    Highlight   `json:"highlight"`
    TextMatch    int64      `json:"text_match"`
}


// Highlight represents the search highlight information
type Highlight struct {
    UserID HighlightDetail `json:"user_id"`
}

// HighlightDetail contains the matched tokens and snippet
type HighlightDetail struct {
    MatchedTokens []string `json:"matched_tokens"`
    Snippet       string   `json:"snippet"`
}

// ParseSearchResponse parses the JSON response from TypeSense into SearchResponse struct
func ParseSearchResponse(responseBody []byte) (*SearchResponse, error) {
    var response SearchResponse
    err := json.Unmarshal(responseBody, &response)
    if err != nil {
        return nil, err
    }
    return &response, nil
}

// ExtractReceipts extracts just the Receipt documents from a SearchResponse
func ExtractProducts(response *SearchResponse) []*model.Product {
    products := make([]*model.Product, len(response.Hits))
    for i, hit := range response.Hits {
        products[i] = &hit.Document
    }
    return products
}

func addProductDocument (httpClient *resty.Client, product model.Product ) error {
	typeSenseURL := os.Getenv("TYPESENSE_URL")
	if typeSenseURL == "" {
		return fmt.Errorf("TYPESENSE_URL is not set")
	}
	resp, err := httpClient.R().
	SetHeader("Content-Type", "application/json").
	SetHeader("X-TYPESENSE-API-KEY", os.Getenv("TYPESENSE_API_KEY")).
	SetBody(map[string]interface{}{
		"id": 		product.ID,
		"user_id": 	product.UserID,
		"unit_price": 	product.UnitPrice,
		"name": 		product.Name,
		"created_at": 	product.CreatedAt,
		"updated_at": 	product.UpdatedAt,
		"deleted_at": 	product.DeletedAt,
		}).
	Post(fmt.Sprintf("%s/collections/products/documents", typeSenseURL ) )

if err != nil {
	return  err
}
if resp.StatusCode() != http.StatusCreated {
	return fmt.Errorf("typesense returned status %d: %s", resp.StatusCode(), resp.String())
}
return  nil
}

func deleteProductDocument (httpClient *resty.Client, id string ) error {
	typeSenseURL := os.Getenv("TYPESENSE_URL")
	if typeSenseURL == "" {
		return fmt.Errorf("TYPESENSE_URL is not set")
	}
	resp, err := httpClient.R().
	SetHeader("X-TYPESENSE-API-KEY", os.Getenv("TYPESENSE_API_KEY")).
	Delete(fmt.Sprintf("%s/collections/products/documents/%s", typeSenseURL, id ) )

	if err != nil {
		return  err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("typesense returned status %d: %s", resp.StatusCode(), resp.String())
	}
	return  nil
}

func searchProductDocuments(httpClient *resty.Client, userId string, query string,) ([]*model.Product, error) {
	typeSenseURL := os.Getenv("TYPESENSE_URL")
	if typeSenseURL == "" {
		return nil, fmt.Errorf("TYPESENSE_URL is not set")
	}
	filters := "user_id:=" + userId // Initialize the filters with user_id condition.
	resp, err := httpClient.R().
		SetHeader("X-TYPESENSE-API-KEY", os.Getenv("TYPESENSE_API_KEY")).
		Get(fmt.Sprintf("%s/collections/products/documents/search?q=%s&query_by=name&page=%d&per_page=10&filter_by=%s", typeSenseURL, query, 1, filters))
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
	
	products := ExtractProducts(response)

	return products, nil
	
	
}