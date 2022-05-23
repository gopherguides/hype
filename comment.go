package hype

import "fmt"

type Comment string

func (tn Comment) Children() Nodes {
	return Nodes{}
}

func (tn Comment) String() string {
	return fmt.Sprintf("<!-- %s -->", string(tn))
}

func (tn Comment) Text() string {
	return string(tn)
}
