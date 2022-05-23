package main

import (
	"context"
	"fmt"
)

// snippet: example
func main() {
	// create a new background context
	ctx := context.Background()

	// call the A function
	// passing in the background context
	A(ctx)
}

func A(ctx context.Context) {
	// wrap the context with a request_id
	// to represent this specific A request
	ctx = context.WithValue(ctx, "request_id", "123")

	// call the B function
	// passing in the wrapped context
	B(ctx)
}

func B(ctx context.Context) {
	// wrap the context with a request_id
	// to represent this specific B request
	ctx = context.WithValue(ctx, "request_id", "456")
	Logger(ctx)
}

// Logger logs the webs request_id
// as well as the request_id from the B
func Logger(ctx context.Context) {
	a := ctx.Value("request_id")
	fmt.Println("A\t", "request_id:", a)

	b := ctx.Value("request_id")
	fmt.Println("B\t", "request_id:", b)
}

// snippet: example

/*
// snippet: out
A	 request_id: 456
B	 request_id: 456
// snippet: out
*/
