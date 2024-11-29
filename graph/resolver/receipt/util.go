package receipt

import (
	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/paginationUtil"
)

func convertEdges(edges []*paginationUtil.Edge[*model.Receipt]) []*model.ReceiptEdge {
    dataEdges := make([]*model.ReceiptEdge, len(edges))
    for i, edge := range edges {
        dataEdges[i] = &model.ReceiptEdge{
            Cursor: edge.Cursor,
            Node:   edge.Node,
        }
    }
    return dataEdges
}

