package hype

import (
	"fmt"

	"github.com/gopherguides/hype/htmx"
	"golang.org/x/net/html"
)

type Text struct {
	*Node
}

func (t Text) StartTag() string {
	return t.String()
}

func (t Text) EndTag() string {
	return ""
}

func (t Text) String() string {
	if t.Node == nil {
		return ""
	}
	return t.Node.Data
	// return fmt.Sprintf("%q", t.Node.Data)
}

func (p *Parser) NewText(node *html.Node) (*Text, error) {
	return NewText(node)
}

func NewText(node *html.Node) (*Text, error) {
	if node == nil {
		return nil, fmt.Errorf("text node can not be nil")
	}

	if node.Type != html.TextNode {
		return nil, fmt.Errorf("node is not a text node %v", node)
	}

	return &Text{
		Node: NewNode(node),
	}, nil
}

func QuickText(s string) *Text {
	n := htmx.TextNode(s)
	nn := NewNode(n)
	return &Text{
		Node: nn,
	}
}
