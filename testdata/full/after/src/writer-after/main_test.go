package main

import (
	"bytes"
	"testing"
)

// snippet: test
func Test_WriteData(t *testing.T) {
	t.Parallel()

	// create a buffer to write to
	bb := &bytes.Buffer{}

	data := []byte("Hello, World!")

	// write the data to the buffer
	WriteData(bb, data)

	// capture the data written to the buffer
	// to the act variable
	act := bb.String()

	exp := string(data)

	// compare the expected and actual values
	if act != exp {
		t.Fatalf("expected %q, got %q", exp, act)
	}

}

// snippet: test
