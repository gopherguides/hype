package main

import (
	"os"
)

func main() {
	WriteData(os.Stdout, []byte("Hello, World!"))
}

func WriteData(w *os.File, data []byte) {
	w.Write(data)
}
