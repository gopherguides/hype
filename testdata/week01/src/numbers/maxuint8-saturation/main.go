package main

import "fmt"

func main() {
	// section: main
	var maxUint8 uint8 = 11
	maxUint8 = maxUint8 * 25
	fmt.Println("new value:", maxUint8)
	// section: main
}

/*
// section: output
new value: 19
// section: output
*/
