package main

import (
	"fmt"
)

func Greet() {
	greet("World!")
}

func greet(s string) {
	fmt.Println("Hello, " + s)
}

func main() {
	Greet()
}
