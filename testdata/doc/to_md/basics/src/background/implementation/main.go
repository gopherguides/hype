package main

import (
	"context"
	"fmt"
)

// snippet: example
func main() {
	ctx := context.Background()

	// print the current value
	// of the context
	fmt.Printf("%v\n", ctx)

	// print Go-syntax representation of the value
	fmt.Printf("\t%#v\n", ctx)

	// print the value of the Done channel
	// does not block because we are not
	// trying to read/write to the channel
	fmt.Printf("\tDone:\t%#v\n", ctx.Done())

	// print the value of the Err
	fmt.Printf("\tErr:\t%#v\n", ctx.Err())

	// print the value of "KEY"
	fmt.Printf("\tValue:\t%#v\n", ctx.Value("KEY"))

	// print the deadline time
	// and true/false if there is no deadline
	deadline, ok := ctx.Deadline()
	fmt.Printf("\tDeadline:\t%s (%t)\n", deadline, ok)
}

// snippet: example

/*
// snippet: out
context.Background
	(*context.emptyCtx)(0xc000016100)
	Done:	(<-chan struct {})(nil)
	Err:	<nil>
	Value:	<nil>
	Deadline:	0001-01-01 00:00:00 +0000 UTC (false)
// snippet: out
*/
