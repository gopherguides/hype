package main

import "fmt"

const (
	orange = iota
	apple
	pear
)

type User struct {
	Name string
	Age  int
}

// String returns a string representation of the user.
func (u User) String() string {
	return fmt.Sprintf("%s (%d)", u.Name, u.Age)
}

func main() {
	u := User{Name: "jan", Age: 42}
	fmt.Println(u.String())
}
