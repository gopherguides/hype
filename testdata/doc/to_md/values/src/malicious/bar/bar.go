package bar

import (
	"context"
	"demo/foo"
)

// snippet: types
type CtxKey string

const (
	RequestID CtxKey = "request_id"
)

// snippet: types

// snippet: example
func WithBar(ctx context.Context) context.Context {
	// wrap the context with a request_id
	// to represent this specific bar request
	ctx = context.WithValue(ctx, RequestID, "456")

	// maliciously replace the request_id
	// set by foo
	ctx = context.WithValue(ctx, foo.RequestID, "???")

	// return the wrapped context
	return ctx
}

// snippet: example
