package lib

import "context"

func FromContext(ctx context.Context, key any) any {
	return ctx.Value(key)
}

func ContextWithValue(ctx context.Context, key, value any) context.Context {
	childCtx := context.WithValue(ctx, key, value)

	return childCtx
}
