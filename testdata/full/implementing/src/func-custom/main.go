package main

import (
	"io"
	"os"
)

// snippet: example
func main() {
	WriteData(os.Stdout, []byte("Hello, World!"))
}

// snippet: example

// snippet: def
func WriteData(w io.Writer, data []byte) {
	w.Write(data)
}

// snippet: def
