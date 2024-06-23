package userLoader

import (
	"context"
	"fmt"

	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"gorm.io/gorm"
)

type UserReader struct {
	db *gorm.DB
}

func NewUserReader(db *gorm.DB) *UserReader {
	return &UserReader{db: db}
}
func (u *UserReader) GetUsers(ctx context.Context, userIds []string) ([]*model.User, []error) {
	if len(userIds) == 0 {
		return []*model.User{}, nil
	}

	var users []*model.User
	err := u.db.Where("id IN (?)", userIds).Find(&users).Error
	if err != nil {
		return nil, []error{fmt.Errorf("failed to fetch users: %w", err)}
	}

	return users, nil
}