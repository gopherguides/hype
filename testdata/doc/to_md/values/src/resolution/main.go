package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/markbates/contextual"
)

// CtxKey is used for the context keys
type CtxKey string

// snippet: example
func main() {

	// create a new background context
	ctx := context.Background()

	// wrap the context with a new context
	// that has the key "A" and the value "a",
	ctx = context.WithValue(ctx, CtxKey("A"), "a")

	// wrap the context with a new context
	// that has the key "B" and the value "b",
	ctx = context.WithValue(ctx, CtxKey("B"), "b")

	// wrap the context with a new context
	// that has the key "C" and the value "c",
	ctx = context.WithValue(ctx, CtxKey("C"), "c")

	// print the final context
	print("ctx", ctx)

	// retreive and print the value
	// for the key "A"
	a := ctx.Value(CtxKey("A"))
	fmt.Println("A:", a)

	// retreive and print the value
	// for the key "B"
	b := ctx.Value(CtxKey("B"))
	fmt.Println("B:", b)

	// retreive and print the value
	// for the key "C"
	c := ctx.Value(CtxKey("C"))
	fmt.Println("C:", c)

}

// snippet: example

func print(label string, ctx context.Context) {
	fmt.Printf("%s.", label)
	err := contextual.Print(ctx, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}
