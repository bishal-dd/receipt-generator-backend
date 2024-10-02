package profile

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/redisUtil"
)

func (r *ProfileResolver) ProfileByUserID(ctx context.Context, userId string) (*model.Profile, error) {
	profile, err := r.GetCachedProfileByUserId(ctx, userId)
	if profile == nil {
		newProfile, err := r.GetProfileByUserID(userId)
		if err != nil {
			return nil, err
		}
		if err = redisUtil.CacheResult(r.redis, ctx, ProfileKey + userId, newProfile, 10); err != nil {
			return nil, err
		}
		return newProfile, nil
	} else if err != nil {
		return nil, err
	}
	
	return profile, nil
}

func (r *ProfileResolver) Profile(ctx context.Context, id string) (*model.Profile, error) {
	var profile *model.Profile
	if err := r.db.Where("id = ?", id).First(&profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}