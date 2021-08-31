package main

import "fmt"

func main() {
	// snippet: main
	var maxUint8 uint8 = 11
	maxUint8 = maxUint8 * 25
	fmt.Println("new value:", maxUint8)
	// snippet: main
}

/*
// snippet: output
new value: 19
// snippet: output
*/
