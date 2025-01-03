package receipt

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/contextUtil"
	"github.com/bishal-dd/receipt-generator-backend/helper/paginationUtil"
	"github.com/bishal-dd/receipt-generator-backend/helper/search"
)


func (r *ReceiptResolver) Receipts(ctx context.Context, first *int, after *string) (*model.ReceiptConnection, error) {
    userId, err := contextUtil.UserIdFromContext(ctx)
    if err != nil {
        return nil, err
    }
    
    offset, limit, err := paginationUtil.CalculatePagination(first, after)
	if err != nil {
		return nil, err 
	} 
    totalReceipts, err := r.CountTotalReceipts()
    if err != nil {
        return nil, err
    }
    receipts, err := r.FetchReceiptsFromDB(ctx, offset, limit, userId)
    if err != nil {
        return nil, err
    }

    connection := paginationUtil.CreateConnection(receipts, totalReceipts, offset)

    return &model.ReceiptConnection{
        Edges: convertEdges(connection.Edges),
        PageInfo: (*model.PageInfo)(connection.PageInfo),
        TotalCount: connection.TotalCount,
    }, nil
}


func (r *ReceiptResolver) Receipt(ctx context.Context, id string) (*model.Receipt, error) {
	newReceipt, err := r.GetReceiptFromDB( ctx, id)
    if err != nil {
		return nil, err
	}
        
    return newReceipt, nil
}

func (r *ReceiptResolver) SearchReceipts(ctx context.Context, page int, year *int, date *string, dateRange []string) (*model.SearchReceipt, error) {
    userId, err := contextUtil.UserIdFromContext(ctx)
    if err != nil {
        return nil, err
    }
    if err := searchDataRangeValidation(dateRange); err != nil {
        return nil, err
    }
    response, err := search.SearchReceiptDocuments(r.httpClient, userId, page, year, date, dateRange)
    if err != nil {
        return nil, err
    }
   
    return response, nil
}