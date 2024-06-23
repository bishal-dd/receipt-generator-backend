package loaders

import "context"

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

