package hype

import (
	"encoding/json"
	"fmt"
)

type Comment string

func (c Comment) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"text": string(c),
		"type": toType(c),
	}
	return json.MarshalIndent(m, "", "  ")
}

func (tn Comment) Children() Nodes {
	return Nodes{}
}

func (tn Comment) String() string {
	return fmt.Sprintf("<!-- %s -->", string(tn))
}

func (tn Comment) Text() string {
	return string(tn)
}
