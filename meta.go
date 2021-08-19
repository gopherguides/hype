package hype

import (
	"fmt"
)

var _ Tag = &Meta{}

// type MetaData map[string]string

type Meta struct {
	*Node
	Key string
	Val string
}

func (m Meta) String() string {
	return m.InlineTag()
}

func (p *Parser) NewMeta(node *Node) (*Meta, error) {
	if node == nil {
		return nil, fmt.Errorf("node can not be nil")
	}

	ats := node.Attrs()

	if ch, ok := ats["charset"]; ok {
		ats["property"] = "charset"
		ats["content"] = ch
	}

	prop, pok := ats["property"]
	name, nok := ats["name"]

	if pok && nok {
		return nil, fmt.Errorf("both property and name defined, pick one %v", node)
	}

	if !pok && !nok {
		return nil, fmt.Errorf("missing property/name %v", node)
	}

	if len(prop) == 0 {
		prop = name
	}

	val, ok := ats["content"]
	if !ok {
		return nil, fmt.Errorf("missing content %v", node)
	}

	m := &Meta{
		Node: node,
		Key:  prop,
		Val:  val,
	}

	for k, v := range ats {
		m.Set(k, v)
	}

	return m, nil
}
