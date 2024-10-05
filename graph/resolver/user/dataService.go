package user

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/loaders"
	"github.com/bishal-dd/receipt-generator-backend/graph/model"
)

func (r *UserResolver) CountTotalUsers() (int64, error) {
    var totalUsers int64
    if err := r.db.Model(&model.User{}).Count(&totalUsers).Error; err != nil {
        return 0, err
    }
    return totalUsers, nil
}

func (r *UserResolver) FetchUsersFromDB(ctx context.Context, offset, limit int) ([]*model.User, error) {
    loaders := loaders.For(ctx)
    var users []*model.User
    if err := r.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
        return nil, err
    }
    userIds := make([]string, len(users))
    for i, user := range users {
        userIds[i] = user.ID
    }

    receiptResults, err := loaders.ReceiptLoader.LoadAll(ctx, userIds)
    if err != nil {
        return nil, err
    }
    for i, receipts := range receiptResults {
        users[i].Receipts = receipts  
    }

    users, err = r.LoadProfileFromUsers(ctx, users)
    if err != nil {
        return nil, err
    }
    
    return users, nil
}