package main

import "fmt"

func main() {
	// section: main
	a := "Hello, 世界"
	for i, c := range a {
		fmt.Printf("%d: %s\n", i, string(c))
	}
	// section: main
}
