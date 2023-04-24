package hype

import (
	"encoding/json"
	"fmt"
)

type Comment string

func (c Comment) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type": fmt.Sprintf("%T", c),
		"text": string(c),
	})
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
