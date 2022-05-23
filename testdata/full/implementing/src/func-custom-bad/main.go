package main

import (
	"io"
)

// snippet: def
func WriteData(w io.Writer, data []byte) {
	w.Write(data)
}

// snippet: def

// snippet: bad
func main() {
	WriteData(42, []byte("Hello, World!"))
}

// snippet: bad
