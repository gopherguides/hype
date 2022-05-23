package main

import (
	"testing"
)

// snippet: scribe
type Scribe struct {
	data []byte
}

func (s *Scribe) Write(p []byte) (int, error) {
	s.data = p
	return len(p), nil
}

// snippet: stringer
func (s Scribe) String() string {
	return string(s.data)
}

// snippet: stringer

// snippet: scribe

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
