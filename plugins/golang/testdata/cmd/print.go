//go:build !sad
// +build !sad

package main

import (
	"fmt"
	"os"
)

func print() {
	fmt.Fprintln(os.Stderr, "STDERR", "Hello, Error!")
	fmt.Fprintln(os.Stdout, "STDOUT", "Hello, World!")
}
