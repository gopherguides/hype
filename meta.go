package hype

import (
	"fmt"

	"golang.org/x/net/html"
)

var _ Tag = &Meta{}

type Metas []*Meta

// Meta represents a <meta> tag.
type Meta struct {
	*Node
}

func (m Meta) String() string {
	return m.InlineTag()
}

// Validate the meta tag
func (m *Meta) Validate(checks ...ValidatorFn) error {
	if m == nil {
		return fmt.Errorf("nil Meta")
	}

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
