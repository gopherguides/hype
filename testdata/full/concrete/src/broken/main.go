package main

import "fmt"

// snippet: musician
type Musician struct {
	Name string
}

func (m Musician) Perform() {
	fmt.Println(m.Name, "is singing")
}

// snippet: musician

// snippet: poet

type Poet struct {
	Name string
}

func (p Poet) Perform() {
	fmt.Println(p.Name, "is reading poetry")
}

// snippet: poet

// snippet: example
func main() {
	m := Musician{Name: "Kurt"}
	PerformAtVenue(m)

	p := Poet{Name: "Janis"}
	PerformAtVenue(p)
}

// snippet: example

// snippet: func
func PerformAtVenue(m Musician) {
	m.Perform()
}

// snippet: func
