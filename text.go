package hype

import (
	"fmt"

	"golang.org/x/net/html"
)

type Text struct {
	*Node
}

func (t Text) String() string {
	if t.Node == nil {
		return ""
	}
	return t.Node.Data
	// return fmt.Sprintf("%q", t.Node.Data)
}

func (p *Parser) NewText(node *html.Node) (*Text, error) {
	if node == nil {
		return nil, fmt.Errorf("node can not be nil")
	}

	if node.Type != html.TextNode {
		return nil, fmt.Errorf("node is not a text node %v", node)
	}

	return &Text{
		Node: NewNode(node),
	}, nil
}
