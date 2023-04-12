package main

import (
	"context"
)

// $ golangci-lint run
func main() {
	// snippet: example
	ctx := context.Background()

	// strings shouldn't be used as keys
	// because they can easily collide
	// with other functions, libraries, etc.
	// that set that same key.
	// instead strings should wrapped in their
	// own type.
	ctx = context.WithValue(ctx, "key", "value")

	// keys must be comparable.
	// maps, and other complex types,
	// are not comparable and can't be used
	// used as keys.
	ctx = context.WithValue(ctx, map[string]int{}, "another value")
	// snippet: example

	// snippet: custom
	// defining a custom string type will
	// help prevent key collisions.
	type CtxKey string
	ctx = context.WithValue(ctx, CtxKey("key"), "value")
	// snippet: custom

	_ = ctx
}
