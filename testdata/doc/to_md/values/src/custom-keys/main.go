package main

import (
	"context"
	"fmt"
)

// snippet: types
// CtxKeyA is used to wrap keys
// associated with a A request
// 	CtxKeyA("request_id")
// 	CtxKeyA("user_id")
type CtxKeyA string

// CtxKeyB is used to wrap keys
// associated with a B request
// 	CtxKeyB("request_id")
// 	CtxKeyB("user_id")
type CtxKeyB string

// snippet: types

func main() {
	// create a new background context
	ctx := context.Background()

	// call A with the background context
	A(ctx)
}

// snippet: example
func A(ctx context.Context) {
	// wrap the context with a request_id
	// to represent this specific A request
	key := CtxKeyA("request_id")
	ctx = context.WithValue(ctx, key, "123")

	// call B with the wrapped context
	B(ctx)
}

func B(ctx context.Context) {
	// wrap the context with a request_id
	// to represent this specific B request
	key := CtxKeyB("request_id")
	ctx = context.WithValue(ctx, key, "456")

	Logger(ctx)
}

// snippet: example

// snippet: logger
// Logger logs the webs request_id
// as well as the request_id from the B
func Logger(ctx context.Context) {
	// retreive the request_id from the A request
	aKey := CtxKeyA("request_id")
	aVal := ctx.Value(aKey)

	// print the request_id from the A request
	print("A", aKey, aVal)

	// retreive the request_id from the B request
	bKey := CtxKeyB("request_id")
	bVal := ctx.Value(bKey)

	// print the request_id from the B request
	print("B", bKey, bVal)
}

// snippet: logger

func print(label string, key any, val any) {
	fmt.Printf("%s: %[2]T(%[2]s): %v\n", label, key, val)
}

/*
// snippet: out
A: main.CtxKeyA(request_id): 123
B: main.CtxKeyB(request_id): 456
// snippet: out
*/
