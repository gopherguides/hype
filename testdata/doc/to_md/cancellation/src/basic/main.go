package main

import (
	"context"
	"fmt"
	"time"
)

// snippet: listener
func listener(ctx context.Context, i int) {
	fmt.Printf("listener %d is waiting\n", i)

	// this will block until the context
	// given context is canceled
	<-ctx.Done()

	fmt.Printf("listener %d is exiting\n", i)
}

// snippet: listener

// snippet: main
func main() {

	// create a background context
	ctx := context.Background()

	// wrap the context with the ability
	// to cancel it
	ctx, cancel := context.WithCancel(ctx)

	// defer cancellation of the context
	// to ensure that any resources are
	// cleaned up regardless of how the
	// function exits
	defer cancel()

	// create 5 listeners
	for i := 0; i < 5; i++ {

		// launch listener in a goroutine
		go listener(ctx, i)

	}

	// allow the listeners to start
	time.Sleep(time.Millisecond * 500)

	fmt.Println("canceling the context")

	// cancel the context and tell the
	// listeners to exit
	cancel()

	// allow the listeners to exit
	time.Sleep(time.Millisecond * 500)
}

// snippet: main

/*
// snippet: out
listener 0 is waiting
listener 3 is waiting
listener 2 is waiting
listener 1 is waiting
listener 4 is waiting
canceling the context
listener 4 is exiting
listener 0 is exiting
listener 1 is exiting
listener 3 is exiting
listener 2 is exiting
// snippet: out
*/
