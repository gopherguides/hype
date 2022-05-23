package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"
)

// snippet: def
func WriteNow(i any) error {
	w, ok := i.(io.Writer)
	if !ok {
		return fmt.Errorf("expected io.Writer, got %T", i)
	}

	now := time.Now()
	w.Write([]byte(now.String()))

	return nil
}

// snippet: def

func main() {
	// snippet: good-assert
	bb := &bytes.Buffer{}
	if err := WriteNow(bb); err != nil {
		log.Fatal(err)
	}
	// snippet: good-assert

	// snippet: bad-assert
	if err := WriteNow(42); err != nil {
		log.Fatal(err)
	}
	// snippet: bad-assert
}

/*
// snippet: panic
expected io.Writer, got int
// snippet: panic
*/
