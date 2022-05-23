package main

import (
	"io"
	"os"
)

func main() {
	WriteData(os.Stdout, []byte("Hello, World!"))
}

func WriteData(w io.Writer, data []byte) {
	w.Write(data)
}
