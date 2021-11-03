package hype

import (
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
}

func (t Text) Validate(checks ...ValidatorFn) error {
	return t.Node.Validate(html.TextNode, checks...)
}

func NewText(node *html.Node) (*Text, error) {
	t := &Text{
		Node: NewNode(node),
	}

	return t, t.Validate()
}

func QuickText(s string) *Text {
	n := htmx.TextNode(s)
	nn := NewNode(n)
	return &Text{
		Node: nn,
	}
}

func (p *Parser) NewText(node *html.Node) (*Text, error) {
	return NewText(node)
}
