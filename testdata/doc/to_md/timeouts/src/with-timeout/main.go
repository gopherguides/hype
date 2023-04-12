package main

import (
	"context"
	"os"
	"time"

	"github.com/markbates/contextual"
)

// snippet: example
func main() {

	// create a background context
	ctx := context.Background()

	// create a new context with a timeout
	// that will cancel the context after 10ms
	// 	equivalent to:
	//		context.WithDeadline(ctx, time.Now().Add(10 *time.Millisecond))
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	print(ctx)
}

// snippet: example

func print(ctx context.Context) {
	// print the context
	contextual.Print(ctx, os.Stdout)
}
