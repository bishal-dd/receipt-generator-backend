package profile

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
)

func (r *ProfileResolver) GetProfileFromDB(id string) (*model.Profile, error) {
	var profile *model.Profile
	if err := r.db.Where("id = ?", id).First(&profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

func (r *ProfileResolver) GetProfileByUserID(userId string) (*model.Profile, error) {
	var profile *model.Profile
	if err := r.db.Where("user_id = ?", userId).First(&profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

func (r *ProfileResolver) DeleteProfileFromDB(ctx context.Context, id string) (error) {
	profile := &model.Profile{
		ID: id,
	}
	if err := r.db.Delete(profile).Error; err != nil {
		return  err
	}
	return  nil
} 