package main

import "fmt"

// section: main
const (
	a = 2
	b = 2
	c = int32(2)
)

func main() {
	fmt.Printf("a = %[1]d (%[1]T)\n", a)
	fmt.Printf("b = %[1]d (%[1]T)\n", b)
	fmt.Printf("c = %[1]d (%[1]T)\n", c)

	fmt.Printf("a*b = %[1]d (%[1]T)\n", a*b)
	fmt.Printf("a*c = %[1]d (%[1]T)\n", a*c)

	d := 4
	e := int32(4)

	fmt.Printf("a*d = %[1]d (%[1]T)\n", a*d)
	fmt.Printf("a*e = %[1]d (%[1]T)\n", a*e)
}

// section: main

/*
// section: output
a = 2 (int)
b = 2 (int)
c = 2 (int32)
a*b = 4 (int)
a*c = 4 (int32)
a*d = 8 (int)
a*e = 8 (int32)
// section: output
*/
