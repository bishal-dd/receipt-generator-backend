package version

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
)

func (r *VersionResolver) VersionByUserID(ctx context.Context, userID string) (*model.Version, error) {
	var version *model.Version
	if err := r.db.Where("user_id = ?", userID).First(&version).Error; err != nil {
		return nil, err
	}
	return version, nil
}

func (r *VersionResolver) Version(ctx context.Context, id string) (*model.Version, error) {
	version, err := r.GetVersionFromDB(id)
	if err != nil {
		return nil, err
	}
	return version, nil
}