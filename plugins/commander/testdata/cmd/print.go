//go:build !sad
// +build !sad

package main

import (
	"fmt"
	"os"
)

func print() {
	fmt.Fprintln(os.Stderr, "stderr->", "Hello, Error!")
	fmt.Fprintln(os.Stdout, "stdout->", "Hello, World!")
}