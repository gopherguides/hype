package main

import "fmt"

func main() {
	// section: main
	const gopher = "Genny"
	gopher = "george"
	fmt.Println(gopher)
	// section: main
}

/*
// section: output
./main.go:8:9: cannot assign to gopher
// section: output
*/
