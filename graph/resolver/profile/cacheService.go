package profile

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper"
	"github.com/redis/go-redis/v9"
)

func (r *ProfileResolver) GetCachedProfileByUserId(ctx context.Context, userId string) (*model.Profile, error) {
	var profile *model.Profile
	cacheKey := ProfileKey + userId
	profileJSON, err := r.redis.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err := helper.Unmarshal([]byte(profileJSON), &profile); err != nil {
		return nil, err
	}	 
    
	 return profile, nil
}