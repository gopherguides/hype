package hype

import (
	"encoding/json"
	"fmt"
)

type UL struct {
	*Element
}

func (ul *UL) MarshalJSON() ([]byte, error) {
	if ul == nil {
		return nil, ErrIsNil("ul")
	}

	ul.RLock()
	defer ul.RUnlock()

	m, err := ul.JSONMap()
	if err != nil {
		return nil, err
	}

	m["type"] = fmt.Sprintf("%T", ul)

	return json.Marshal(m)
}

func (ol *UL) MD() string {
	if ol == nil || ol.Element == nil {
		return ""
	}

	return ol.Children().MD()
}

func NewULNodes(p *Parser, el *Element) (Nodes, error) {
	if el == nil {
		return nil, nil
	}

	ol := &UL{
		Element: el,
	}

	return Nodes{ol}, nil
}
