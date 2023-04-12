package main

import (
	"context"
	"fmt"
	"time"
)

// snippet: example
func main() {

	// create a background context
	ctx := context.Background()

	// wrap the context that will
	// self cancel after 10 milliseconds
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	// check the error:
	//	<nil>
	fmt.Println("ctx.Err()", ctx.Err())

	// wait for the context to self cancel
	<-ctx.Done()

	// check the error:
	//	context.Canceled
	fmt.Println("ctx.Err()", ctx.Err())

	// check the error again:
	//	context.DeadlineExceeded
	fmt.Println("ctx.Err()", ctx.Err())
}

// snippet: example
