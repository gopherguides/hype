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

func (p *Parser) NewDocType(node *html.Node) (*DocType, error) {
	if node == nil {
		return nil, fmt.Errorf("node can not be nil")
	}

	if node.Type != html.DoctypeNode {
		return nil, fmt.Errorf("node is not a doctype node %v", node)
	}

	return &DocType{
		Node: NewNode(node),
	}, nil
}
