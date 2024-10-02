package version

import (
	"context"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper"
	"github.com/bishal-dd/receipt-generator-backend/helper/redisUtil"
)

func (r *VersionResolver) CreateVersion(ctx context.Context, input model.CreateVersion) (*model.Version, error) {
	newVersion := &model.Version{
		ID: helper.UUID(),
		Mode: input.Mode,
		UserID: input.UserID,
		CreatedAt: helper.CurrentTime(),
	}
	if err := r.db.Create(newVersion).Error; err != nil {
		return nil, err
	}
	return newVersion, nil
}

func (r *VersionResolver) UpdateVersion(ctx context.Context, input model.UpdateVersion) (*model.Version, error) {
	version := &model.Version{
		ID: input.ID,
	}
	
	if err := r.db.Model(version).Updates(input).Error; err != nil {
		return nil, err
	}
	if err := redisUtil.DeleteCacheItem(r.redis, ctx, VersionKey, input.ID); err != nil {
		return nil, err
	}
	newVersion, err := r.GetVersionFromDB(input.ID)
	if err != nil {
		return nil, err
	}

	return newVersion, nil
}

func (r *VersionResolver) DeleteVersion(ctx context.Context, id string) (bool, error) {
	if err := r.DeleteVersionFromDB(ctx, id); err != nil {
		return false, err
	}
	if err := redisUtil.DeleteCacheItem(r.redis, ctx, VersionKey, id); err != nil {
		return false, err
	}
	
	return true, nil
}