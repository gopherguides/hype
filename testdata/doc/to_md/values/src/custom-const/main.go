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

// snippet: consts
const (
	// A_RequestID can be used to
	// retreive the request_id for
	// the A request
	A_RequestID CtxKeyA = "request_id"
	// 	A_SESSION_ID CtxKeyA = "session_id"
	// 	A_SERVER_ID CtxKeyA = "server_id"
	// 	other keys...

	// B_RequestID can be used to
	// retreive the request_id for
	// the B request
	B_RequestID CtxKeyB = "request_id"
)

// snippet: consts

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
	ctx = context.WithValue(ctx, A_RequestID, "123")

	// call B with the wrapped context
	B(ctx)
}

func B(ctx context.Context) {
	// wrap the context with a request_id
	// to represent this specific B request
	ctx = context.WithValue(ctx, B_RequestID, "456")

	Logger(ctx)
}

// snippet: example

// snippet: logger
// Logger logs the webs request_id
// as well as the request_id from the B
func Logger(ctx context.Context) {
	// retreive the request_id from the A request
	aKey := A_RequestID
	aVal := ctx.Value(aKey)

	// print the request_id from the A request
	print("A", aKey, aVal)

	// retreive the request_id from the B request
	bKey := B_RequestID
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
