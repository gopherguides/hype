package foo

import (
	"context"
	"fmt"
)

// snippet: types
type ctxKey string

const (
	requestID ctxKey = "request_id"
)

// snippet: types

func WithFoo(ctx context.Context) context.Context {
	// wrap the context with a request_id
	// to represent this specific foo request
	ctx = context.WithValue(ctx, requestID, "123")

	// return the wrapped context
	return ctx
}

// snippet: example
func RequestIDFrom(ctx context.Context) (string, error) {
	// get the request_id from the context
	s, ok := ctx.Value(requestID).(string)
	if !ok {
		return "", fmt.Errorf("request_id not found in context")
	}
	return s, nil
}

// snippet: example
