package main

import "fmt"

const (
	year     = 365        // untyped
	leapYear = int32(366) // typed
)

func main() {
	hours := 24
	minutes := int32(60)
	fmt.Println(hours * year)       // multiplying an int and untyped
	fmt.Println(minutes * year)     // multiplying an int32 and untyped
	fmt.Println(minutes * leapYear) // multiplying both int32 types
}

/*
	output:
	8760
	21900
	21960
*/
