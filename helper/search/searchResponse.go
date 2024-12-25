package search

import (
	"encoding/json"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
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
    Document     model.Receipt     `json:"document"`
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
func ExtractReceipts(response *SearchResponse) []*model.Receipt {
    receipts := make([]*model.Receipt, len(response.Hits))
    for i, hit := range response.Hits {
        receipts[i] = &hit.Document
    }
    return receipts
}