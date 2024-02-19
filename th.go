package hype

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type TH struct {
	*Element
}

func (th *TH) MarshalJSON() ([]byte, error) {
	if th == nil {
		return nil, ErrIsNil("th")
	}

	th.RLock()
	defer th.RUnlock()

	m, err := th.JSONMap()
	if err != nil {
		return nil, err
	}

	m["type"] = fmt.Sprintf("%T", th)

	return json.MarshalIndent(m, "", "  ")
}

func (th *TH) IsEmptyNode() bool {
	if th == nil {
		return true
	}

	kids := th.Children()
	if len(kids) == 0 {
		return true
	}

	return IsEmptyNode(kids)
}

func NewTH(p *Parser, el *Element) (*TH, error) {
	if el == nil {
		return nil, ErrIsNil("th")
	}

	th := &TH{
		Element: el,
	}

	body := th.Children().String()

	if len(body) == 0 {
		return th, nil
	}

	nodes, err := p.ParseFragment(strings.NewReader(body))
	if err != nil {
		if !errors.Is(err, ErrNilFigure) {
			return nil, th.WrapErr(err)
		}
	}

	th.Nodes = nodes

	return th, nil
}

func NewTHNodes(p *Parser, el *Element) (Nodes, error) {
	th, err := NewTH(p, el)
	if err != nil {
		return nil, err
	}

	return Nodes{th}, nil
}
