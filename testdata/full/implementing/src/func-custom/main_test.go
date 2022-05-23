package main

import (
	"testing"
)

// snippet: test
func Test_WriteData(t *testing.T) {
	t.Parallel()

	scribe := &Scribe{}
	data := []byte("Hello, World!")
	WriteData(scribe, data)

	act := scribe.String()
	exp := string(data)
	if act != exp {
		t.Fatalf("expected %q, got %q", exp, act)
	}

}

// snippet: test
