package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("bad exit status")
	os.Exit(1)
}
