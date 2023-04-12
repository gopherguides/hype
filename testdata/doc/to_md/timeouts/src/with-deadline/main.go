package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/markbates/contextual"
)

// snippet: example
func main() {

	// create a background context
	ctx := context.Background()

	// create an absolute date/time (January 1, 2030)
	deadline := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	fmt.Println("deadline:", deadline.Format(time.RFC3339))

	// create a new context with a deadline
	// that will cancel at January 1, 2030 00:00:00.
	ctx, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()

	print(ctx)
}

// snippet: example

func print(ctx context.Context) {
	// print the context
	contextual.Print(ctx, os.Stdout)
}

/*
// snippet: out
deadline: 2030-01-01T00:00:00Z
WithTimeout(deadline: {wall:0 ext:64029052800 loc:<nil>})
	--> Background
// snippet: out
*/
