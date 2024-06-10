package receipt

import (
	"encoding/base64"
	"strconv"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
)

func Edges(offset int, limit int, values []*model.Receipt,) ([]*model.ReceiptEdge, int) {
	end := offset + limit
	if end > len(values) {
		end = len(values)
	}
	// Safely create edges
	edges := make([]*model.ReceiptEdge, end-offset)
	for i := offset; i < end; i++ {
		cursor := "cursor" + strconv.Itoa(i)
		encodedCursor := base64.StdEncoding.EncodeToString([]byte(cursor))
		edges[i-offset] = &model.ReceiptEdge{
			Cursor: encodedCursor,
			Node:   values[i],
		}
	}

	return edges, end
}

func PageInfo (edges []*model.ReceiptEdge, totalReceipts int64, end int, offset int) *model.PageInfo {
    var startCursor, endCursor *string
    if len(edges) > 0 {
        startCursor = &edges[0].Cursor
        endCursor = &edges[len(edges)-1].Cursor
    }

    pageInfo := &model.PageInfo{
        HasNextPage:     int64(end) < totalReceipts,
        HasPreviousPage: offset > 0,
        StartCursor:     startCursor,
        EndCursor:       endCursor,
    }
	return pageInfo
}