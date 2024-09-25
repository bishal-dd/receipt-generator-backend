package contextUtil

import (
	"context"
)

type Key string

func SetContextValue(ctx context.Context, key Key, value interface{}) context.Context {
	return context.WithValue(ctx, key, value)
}