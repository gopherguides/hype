package foo

import "context"

// snippet: types
type CtxKey string

const (
	RequestID CtxKey = "request_id"
)

// snippet: types

// snippet: example
func WithFoo(ctx context.Context) context.Context {
	// wrap the context with a request_id
	// to represent this specific foo request
	ctx = context.WithValue(ctx, RequestID, "123")

	// return the wrapped context
	return ctx
}

// snippet: example
