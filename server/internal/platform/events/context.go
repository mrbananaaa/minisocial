package events

import "context"

type contextKey struct{}

func ContextWithCollector(
	ctx context.Context,
	collector *Collector,
) context.Context {
	return context.WithValue(ctx, contextKey{}, collector)
}

func CollectorFromContext(
	ctx context.Context,
) (*Collector, bool) {
	c, ok := ctx.Value(contextKey{}).(*Collector)
	return c, ok
}
