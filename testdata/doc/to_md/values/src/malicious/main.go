package main

import (
	"context"
	"demo/bar"
	"demo/foo"
	"fmt"
)

// snippet: example
func main() {
	// create a background context
	ctx := context.Background()

	// wrap the context with foo
	ctx = foo.WithFoo(ctx)

	// wrap the context with bar
	ctx = bar.WithBar(ctx)

	// retrieve the foo.RequestID
	// value from the context
	id := ctx.Value(foo.RequestID)

	// print the value
	fmt.Println("foo.RequestID: ", id)
}

// snippet: example

/*
// snippet: out
foo.RequestID:  ???
// snippet: out
*/
