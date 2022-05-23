package demo

import (
	"testing"
)

// snippet: example
func Test_Slicer(t *testing.T) {
	t.Parallel()

	input := "Hello World"

	act := Slicer(input)

	exp := []string{input}

	if len(act) != len(exp) {
		t.Fatalf("expected %v, got %v", exp, act)
	}

	for i, v := range exp {
		if act[i] != v {
			t.Fatalf("expected %v, got %v", exp, act)
		}
	}

}

// snippet: example
