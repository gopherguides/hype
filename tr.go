package hype

import (
	"encoding/json"
	"fmt"
)

type TR struct {
	*Element
}

func (tr *TR) MarshalJSON() ([]byte, error) {
	if tr == nil {
		return nil, ErrIsNil("tr")
	}

	m, err := tr.JSONMap()
	if err != nil {
		return nil, err
	}

	m["type"] = fmt.Sprintf("%T", tr)

	return json.MarshalIndent(m, "", "  ")
}

func (tr *TR) IsEmptyNode() bool {
	if tr == nil {
		return true
	}

	kids := tr.Children()
	if len(kids) == 0 {
		return true
	}

	return IsEmptyNode(kids)
}

func NewTR(el *Element) (*TR, error) {
	if el == nil {
		return nil, ErrIsNil("tr")
	}

	tr := &TR{
		Element: el,
	}

	return tr, nil
}

func NewTRNodes(p *Parser, el *Element) (Nodes, error) {
	tr, err := NewTR(el)
	if err != nil {
		return nil, err
	}

	return Nodes{tr}, nil
}
