package bar

import (
	"context"
)

// snippet: types
type ctxKey string

const (
	requestID ctxKey = "request_id"
)

// snippet: types

// snippet: example
func WithBar(ctx context.Context) context.Context {
	// wrap the context with a request_id
	// to represent this specific bar request
	ctx = context.WithValue(ctx, requestID, "456")

	// no longer able to set the foo request id
	// it does not have access to the foo.ctxKey type
	// as it is not exported, so bar can not create
	// a new key of that type.
	// ctx = context.WithValue(ctx, foo.ctxKey("request_id"), "???")

	// return the wrapped context
	return ctx
}

// snippet: example
