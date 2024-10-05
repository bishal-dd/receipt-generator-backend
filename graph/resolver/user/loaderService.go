package user

import (
	"context"
	"fmt"

	"github.com/bishal-dd/receipt-generator-backend/graph/loaders"
	"github.com/bishal-dd/receipt-generator-backend/graph/model"
)


func (r *UserResolver) LoadProfileFromUsers(ctx context.Context, users []*model.User) ([]*model.User, error) {
	loaders := loaders.For(ctx)
	userIds := make([]string, len(users))
	for i, user := range users {
		userIds[i] = user.ID
	}

	profileResults, err := loaders.ProfileLoader.LoadAll(ctx, userIds)
	if err != nil {
		return nil, fmt.Errorf("failed to load profiles: %w", err)
	}

	fmt.Printf("Loaded profiles: %+v\n", profileResults)

	for i, profile := range profileResults {
		if profile != nil {
			users[i].Profile = profile[0]
		} else {
			fmt.Printf("No profile found for user: %s\n", users[i].ID)
		}
	}

	return users, nil
}
