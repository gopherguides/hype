package hype

import (
	"encoding/json"
	"errors"
	"strings"
)

type TD struct {
	*Element
}

func (td *TD) MarshalJSON() ([]byte, error) {
	if td == nil {
		return nil, ErrIsNil("td")
	}

	td.RLock()
	defer td.RUnlock()

	m, err := td.JSONMap()
	if err != nil {
		return nil, err
	}

	m["type"] = toType(td)

	return json.MarshalIndent(m, "", "  ")
}

func (td *TD) IsEmptyNode() bool {
	if td == nil {
		return true
	}

	kids := td.Children()
	if len(kids) == 0 {
		return true
	}

	return IsEmptyNode(kids)
}

func NewTD(p *Parser, el *Element) (*TD, error) {
	if el == nil {
		return nil, ErrIsNil("td")
	}

	td := &TD{
		Element: el,
	}

	body := td.Children().String()

	if len(body) == 0 {
		return td, nil
	}

	nodes, err := p.ParseFragment(strings.NewReader(body))
	if err != nil {
		if !errors.Is(err, ErrNilFigure) {
			return nil, td.WrapErr(err)
		}
	}

	td.Nodes = nodes

	return td, nil
}

func NewTDNodes(p *Parser, el *Element) (Nodes, error) {
	td, err := NewTD(p, el)
	if err != nil {
		return nil, err
	}

	return Nodes{td}, nil
}
