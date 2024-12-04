package profile

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/cloudFront"
)

func (r *ProfileResolver) ProfileByUserID(ctx context.Context, userId string) (*model.Profile, error) {
	profile, err := r.GetProfileByUserID(userId)
	if err != nil {
		return nil, err
	}

	if profile.SignatureImage != nil && *profile.SignatureImage != "" {
		signedURL, err := cloudFront.GetCloudFrontURL(*profile.SignatureImage)
		if err != nil {
			return nil, err
		}
		profile.SignatureImage = &signedURL
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