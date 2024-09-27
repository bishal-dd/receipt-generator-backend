package user

import (
	"context"
	"fmt"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper"
	"github.com/redis/go-redis/v9"
)

func (r *UserResolver) GetCachedUsers(ctx context.Context, userId string, offset, limit int) ([]*model.User, error) {
    pageCacheKey := fmt.Sprintf("%s:%d:%d:%s", UsersKey, offset, limit, userId)
    usersJSON, err := r.redis.Get(ctx, pageCacheKey).Result()
    if err == redis.Nil {
        return nil, nil
    } else if err != nil {
        return nil, err
    }

    var users []*model.User
    if err := helper.Unmarshal([]byte(usersJSON), &users); err != nil {
        return nil, err
    }
    return users, nil
}