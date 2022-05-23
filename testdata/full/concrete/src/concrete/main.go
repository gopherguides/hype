package main

import "fmt"

type Musician struct {
	Name string
}

func (m Musician) Perform() {
	fmt.Println(m.Name, "is singing")
}

// snippet: example
func main() {
	m := Musician{Name: "Kurt"}
	PerformAtVenue(m)
}

// snippet: example

func PerformAtVenue(m Musician) {
	m.Perform()
}
