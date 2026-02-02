package contextutils

import (
	"context"
)

type ContextKey string

const (
	ContextKeyReqFrom     ContextKey = "reqFrom"
	ContextKeyReqTracking ContextKey = "reqTracking"
)

func GetContextValue(
	ctx context.Context,
	key ContextKey,
) any {
	out := ctx.Value(key)
	if out == nil {
		out = ""
	}
	return out
}

func SetContextValue(
	ctx context.Context,
	key ContextKey,
	val any,
) context.Context {
	return context.WithValue(
		ctx, key, val,
	)
}
