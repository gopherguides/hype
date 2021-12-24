package hype

import (
	"fmt"

	"golang.org/x/net/html"
)

var _ Tag = &DocType{}

// DocType represents the doctype of a document.
type DocType struct {
	*Node
}

func (dt DocType) String() string {
	return fmt.Sprintf("<!doctype %s>\n", dt.Atom())
}

func (dt DocType) Validate(checks ...ValidatorFn) error {
	return dt.Node.Validate(html.DoctypeNode, checks...)
}

// NewDocType returns a new DocType from the given node.
func (p *Parser) NewDocType(node *html.Node) (*DocType, error) {
	return NewDocType(node)
}

// NewDocType returns a new DocType from the given node.
func NewDocType(n *html.Node) (*DocType, error) {

	dt := &DocType{
		Node: NewNode(n),
	}

	return dt, dt.Validate()
}
