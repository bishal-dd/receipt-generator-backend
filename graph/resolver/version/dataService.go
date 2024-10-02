package version

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
)


func (r *VersionResolver) GetVersionFromDB(id string) (*model.Version, error) {
	var version *model.Version
	if err := r.db.Where("id = ?", id).First(&version).Error; err != nil {
		return nil, err
	}
	return version, nil
}

func (r *VersionResolver) DeleteVersionFromDB(ctx context.Context, id string) (error) {
	version := &model.Version{
		ID: id,
	}
	if err := r.db.Delete(version).Error; err != nil {
		return  err
	}
	return  nil
} 