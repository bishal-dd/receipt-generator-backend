package user

import (
	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/paginationUtil"
)

func convertEdges(edges []*paginationUtil.Edge[*model.User]) []*model.UserEdge {
    dataEdges := make([]*model.UserEdge, len(edges))
    for i, edge := range edges {
        dataEdges[i] = &model.UserEdge{
            Cursor: edge.Cursor,
            Node:   edge.Node,
        }
    }
    return dataEdges
}

