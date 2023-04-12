package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/markbates/contextual"
)

type CtxKey string

// snippet: example
func main() {
	ctx := context.Background()
	A(ctx)
}

// snippet: a
func A(ctx context.Context) {
	ctx = context.WithValue(ctx, CtxKey("name"), "Amy")
	B(ctx)
}

// snippet: a

func B(ctx context.Context) {
	print("B", ctx)
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
