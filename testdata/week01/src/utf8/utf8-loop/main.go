package main

import "fmt"

func main() {
	// section: main
	a := "Hello, 世界" // 9 characters (including the space and comma)
	for i := 0; i < len(a); i++ {
		fmt.Printf("%d: %s\n", i, string(a[i]))
	}

	// section: main
}
