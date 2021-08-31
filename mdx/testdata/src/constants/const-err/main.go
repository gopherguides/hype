package main

import "fmt"

func main() {
	// snippet: main
	const gopher = "Genny"
	gopher = "george"
	fmt.Println(gopher)
	// snippet: main
}

/*
// snippet: output
./main.go:8:9: cannot assign to gopher
// snippet: output
*/
