package main

import "fmt"

const (
	leapYear = int32(366) // typed
)

func main() {
	hours := 24
	fmt.Println(hours * leapYear) // multiplying int and int32 types}
}

/*
	output:
	invalid operation: hours * leapYear (mismatched types int and int32)
*/
