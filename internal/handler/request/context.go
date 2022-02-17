package request

import (
	"context"
)

type key int

const (
	userKey key = iota
	clientIPKey
)

func WithClientIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, clientIPKey, ip)
}

func ClientIPFrom(ctx context.Context) string {
	return ctx.Value(clientIPKey).(string)
}
