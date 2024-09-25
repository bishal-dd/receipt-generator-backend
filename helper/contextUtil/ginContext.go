package contextUtil

import (
	"context"
	"fmt"
)

func GinContextFromContext(ctx context.Context) (context.Context, error) {
	ginContext := ctx.Value(GinContextKey)
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(context.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}