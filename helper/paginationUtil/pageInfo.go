package paginationUtil

import (
	"encoding/base64"
	"strconv"
)

// CreatePageInfo creates pagination information
func CreatePageInfo(edgeCount int, totalCount int64, fetchedCount int, offset int) *PageInfo {
    var startCursor, endCursor *string

    if edgeCount > 0 {
        start := "cursor" + strconv.Itoa(offset)
        end := "cursor" + strconv.Itoa(offset+edgeCount-1)
        startEncoded := base64.StdEncoding.EncodeToString([]byte(start))
        endEncoded := base64.StdEncoding.EncodeToString([]byte(end))
        startCursor = &startEncoded
        endCursor = &endEncoded
    }

    return &PageInfo{
        HasNextPage:     int64(offset+fetchedCount) < totalCount,
        HasPreviousPage: offset > 0,
        StartCursor:     startCursor,
        EndCursor:       endCursor,
    }
}

