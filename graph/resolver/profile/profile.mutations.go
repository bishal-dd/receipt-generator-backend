package profile

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/database"
	"github.com/bishal-dd/receipt-generator-backend/helper/redisUtil"
)


func (r *ProfileResolver) CreateProfile(ctx context.Context, input model.CreateProfile) (*model.Profile, error) {
	inputData := database.CreateFields[model.Profile](input);
	if err := r.db.Create(inputData).Error; err != nil {
		return nil, err
	}

	return inputData, nil
}

func (r *ProfileResolver) UpdateProfile(ctx context.Context, input model.UpdateProfile) (*model.Profile, error) {
	profile := &model.Profile{
		ID: input.ID,
	}
	
	if err := r.db.Model(profile).Updates(input).Error; err != nil {
		return nil, err
	}
	if err := redisUtil.DeleteCacheItem(r.redis, ctx, ProfileKey, input.ID); err != nil {
		return nil, err
	}
	newProfile, err := r.GetProfileFromDB(input.ID)
	if err != nil {
		return nil, err
	}

	return newProfile, nil
}


func (r *ProfileResolver) DeleteProfile(ctx context.Context, id string) (bool, error) {
	if err := r.DeleteProfileFromDB(ctx, id); err != nil {
		return false, err
	}
	if err := redisUtil.DeleteCacheItem(r.redis, ctx, ProfileKey, id); err != nil {
		return false, err
	}
	
	return true, nil
}