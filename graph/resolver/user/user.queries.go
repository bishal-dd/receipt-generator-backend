package user

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/loaders"
	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/paginationUtil"
)


func (r *UserResolver) Users(ctx context.Context, first *int, after *string) (*model.UserConnection, error) {
    offset, limit, err := paginationUtil.CalculatePagination(first, after)
	if err != nil {
		return nil, err 
	} 
    totalUsers, err := r.CountTotalUsers()
    if err != nil {
        return nil, err
    }
    users, err := r.FetchUsersFromDB(ctx, offset, limit)
        if err != nil {
            return nil, err
        }
    
    connection := paginationUtil.CreateConnection(users, totalUsers, offset)

    return &model.UserConnection{
        Edges: convertEdges(connection.Edges),
        PageInfo: (*model.PageInfo)(connection.PageInfo),
        TotalCount: connection.TotalCount,
    }, nil
}

func (r *UserResolver) User(ctx context.Context, id string) (*model.User, error) {
	loaders := loaders.For(ctx)
    user, err := loaders.UserLoader.Load(ctx, id)
    if err != nil {
        return nil, err
    }
    return user, nil
}