package cstx

import "context"

type contextKey string

const ctxKey contextKey = "cstx-id"

func WithCstxID(ctx context.Context, cstxID string) context.Context {
	return context.WithValue(ctx, ctxKey, cstxID)
}

func GetCstxID(ctx context.Context) string {
	val, _ := ctx.Value(ctxKey).(string)
	return val
}
