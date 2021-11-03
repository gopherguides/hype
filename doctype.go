package hype

import (
	"fmt"

	"golang.org/x/net/html"
)

var _ Tag = &DocType{}

type DocType struct {
	*Node
}

func (dt DocType) String() string {
	return fmt.Sprintf("<!doctype %s>\n", dt.Data)
}

func (dt DocType) Validate(checks ...ValidatorFn) error {
	return dt.Node.Validate(html.DoctypeNode, checks...)
}

func (p *Parser) NewDocType(node *html.Node) (*DocType, error) {
	return NewDocType(node)
}

func NewDocType(n *html.Node) (*DocType, error) {

	dt := &DocType{
		Node: NewNode(n),
	}

	return dt, dt.Validate()
}
