package main

import "fmt"

type Musician struct {
	Name string
}

// Perform implicitly implements the Performer interface
func (m Musician) Perform() {
	fmt.Println(m.Name, "is singing")
}

type Performer interface {
	Perform()
}

// snippet: example
func PerformAtVenue(p Performer) {
	p.Perform()
}

func main() {
	m := Musician{Name: "Kurt"}
	PerformAtVenue(m)
}

// snippet: example
