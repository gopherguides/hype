//go:build sad
// +build sad

package main

import (
	"fmt"
	"os"
)

func print() {
	fmt.Fprintln(os.Stderr, "boom!")
	os.Exit(-1)
}
