package hype

import (
	"encoding/json"
	"fmt"
)

type OL struct {
	*Element
}

func (ol *OL) MarshalJSON() ([]byte, error) {
	if ol == nil {
		return nil, ErrIsNil("ol")
	}

	ol.RLock()
	defer ol.RUnlock()

	m, err := ol.JSONMap()
	if err != nil {
		return nil, err
	}

	m["type"] = fmt.Sprintf("%T", ol)

	return json.MarshalIndent(m, "", "  ")
}

func (ol *OL) MD() string {
	if ol == nil || ol.Element == nil {
		return ""
	}

	return ol.Children().MD()
}

func NewOLNodes(p *Parser, el *Element) (Nodes, error) {
	if el == nil {
		return nil, nil
	}

	ol := &OL{
		Element: el,
	}

	return Nodes{ol}, nil
}
