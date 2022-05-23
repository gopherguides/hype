package main

import "fmt"

type Musician struct {
	Name string
}

func (m Musician) Perform() {
	fmt.Println(m.Name, "is singing")
}

type Poet struct {
	Name string
}

func (p Poet) Perform() {
	fmt.Println(p.Name, "is reading poetry")
}

type Performer interface {
	Perform()
}

// snippet: example
func main() {
	m := Musician{Name: "Kurt"}
	PerformAtVenue(m)

	p := Poet{Name: "Janis"}
	PerformAtVenue(p)
}

// snippet: example

func PerformAtVenue(p Performer) {
	p.Perform()
}
