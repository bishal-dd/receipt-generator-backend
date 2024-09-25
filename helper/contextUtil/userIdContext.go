package contextUtil

import (
	"context"
	"fmt"
)

func UserIdFromContext(ctx context.Context) (string, error) {
	userId := ctx.Value(UserIDKey)
	if userId == nil {
		err := fmt.Errorf("could not retrieve userID")
		return "", err
	}
	return userId.(string), nil
}