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
}

// snippet: example

/*
// snippet: out
context.Background
	(*context.emptyCtx)(0xc000016100)
// snippet: out
*/
