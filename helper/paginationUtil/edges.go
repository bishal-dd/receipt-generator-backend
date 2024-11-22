package paginationUtil

import (
	"encoding/base64"
	"strconv"
)

// Edge represents a generic edge in a connection
type Edge[T any] struct {
    Cursor string
    Node   T
}

// Connection represents a generic paginated connection
type Connection[T any] struct {
    Edges      []*Edge[T]
    PageInfo   *PageInfo
    TotalCount int
}

// PageInfo contains pagination metadata
type PageInfo struct {
    HasNextPage     bool
    HasPreviousPage bool
    StartCursor     *string
    EndCursor       *string
}

// CreateEdges creates a slice of edges from a slice of nodes
func CreateEdges[T any](nodes []T, offset int) []*Edge[T] {
    edges := make([]*Edge[T], len(nodes))
    for i := 0; i < len(nodes); i++ {
        cursor := "cursor" + strconv.Itoa(offset+i)
        encodedCursor := base64.StdEncoding.EncodeToString([]byte(cursor))
        edges[i] = &Edge[T]{
            Cursor: encodedCursor,
            Node:   nodes[i],
        }
    }
    return edges
}

// CreateConnection creates a complete connection with edges and page info
func CreateConnection[T any](nodes []T, totalCount int64, offset int) *Connection[T] {
    edges := CreateEdges(nodes, offset)
    pageInfo := CreatePageInfo(len(edges), totalCount, len(nodes), offset)

    return &Connection[T]{
        Edges:      edges,
        PageInfo:   pageInfo,
        TotalCount: int(totalCount),
    }
}