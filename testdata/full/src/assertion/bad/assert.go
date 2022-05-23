package main

import (
	"bytes"
	"io"
	"time"
)

// snippet: def
func WriteNow(i any) {
	w := i.(io.Writer)
	now := time.Now()
	w.Write([]byte(now.String()))
}

// snippet: def

func main() {
	// snippet: good-assert
	bb := &bytes.Buffer{}
	WriteNow(bb)
	// snippet: good-assert

	// snippet: bad-assert
	WriteNow(42)
	// snippet: bad-assert
}
