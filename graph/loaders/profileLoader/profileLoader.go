package profileLoader

import (
	"context"
	"fmt"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"gorm.io/gorm"
)

type ProfileReader struct {
	db *gorm.DB
}

func NewProfileReader(db *gorm.DB) *ProfileReader {
	return &ProfileReader{db: db}
}
func (u *ProfileReader) GetProfileByUserIds(ctx context.Context, userIds []string) ([][]*model.Profile, []error) {
    if len(userIds) == 0 {
        return [][]*model.Profile{}, nil
    }

    var profiles []*model.Profile
    err := u.db.Where("user_id IN (?)", userIds).Find(&profiles).Error
    if err != nil {
        return nil, []error{fmt.Errorf("failed to fetch profiles: %w", err)}
    }

    profileMap := make(map[string][]*model.Profile)
    for _, profile := range profiles {
        profileMap[profile.UserID] = append(profileMap[profile.UserID], profile)
    }

    result := make([][]*model.Profile, len(userIds))
    for i, userID := range userIds {
        result[i] = profileMap[userID] 
    }

    return result, nil
}
