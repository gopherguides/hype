package hype

import (
	"fmt"

	"golang.org/x/net/html"
)

var _ Tag = &Meta{}

type Metas []*Meta

// Value returns the value for the key in the <meta> tags.
func (ms Metas) Value(key string) (string, bool) {
	for _, m := range ms {
		if m.Key == key {
			return m.Val, true
		}
	}
	return "", false
}

// Meta represents a <meta> tag.
type Meta struct {
	*Node
	Key string
	Val string
}

func (m Meta) String() string {
	return m.InlineTag()
}

// Validate the meta tag
func (m *Meta) Validate(checks ...ValidatorFn) error {
	if m == nil {
		return fmt.Errorf("nil Meta")
	}

	fn := func(node *Node) error {
		ats := node.Attrs()

		if ch, ok := ats["charset"]; ok {
			ats["property"] = "charset"
			ats["content"] = ch
		}

		if len(m.Key) > 0 && len(m.Val) > 0 {
			return nil
		}

		prop, pok := ats["property"]
		name, nok := ats["name"]

		if pok && nok {
			return fmt.Errorf("both property and name defined, pick one %v", node)
		}

		if !pok && !nok {
			return fmt.Errorf("missing property/name %v", node)
		}

		if len(prop) == 0 {
			prop = name
		}

		val, ok := ats["content"]
		if !ok {
			return fmt.Errorf("missing content %v", node)
		}
		m.Key = prop
		m.Val = val
		return nil
	}

	checks = append(checks, fn)
	return m.Node.Validate(html.ElementNode, checks...)
}

// MetaNode returns a meta tag from the given node.
func NewMeta(node *Node) (*Meta, error) {
	m := &Meta{
		Node: node,
	}

	if err := m.Validate(); err != nil {
		return nil, err
	}

	return m, m.Validate()
}
