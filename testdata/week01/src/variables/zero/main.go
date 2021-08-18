package main

import "fmt"

func main() {
	var a int
	var b string
	var c float64
	var d bool

	fmt.Printf("var a %T = %+v\n", a, a)
	fmt.Printf("var b %T = %q\n", b, b)
	fmt.Printf("var c %T = %+v\n", c, c)
	fmt.Printf("var d %T = %+v\n\n", d, d)
}

/* output
var a int =  0
var b string = ""
var c float64 = 0
var d bool = false
*/
