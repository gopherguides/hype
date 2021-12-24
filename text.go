package hype

import (
	"github.com/gopherguides/hype/htmx"
	"golang.org/x/net/html"
)

// Text represents a text node in the HTML document.
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
	return t.html.Data
}

func (t Text) Validate(checks ...ValidatorFn) error {
	return t.Node.Validate(html.TextNode, checks...)
}

// NewText creates a new Text from an html.Node.
func NewText(node *html.Node) (*Text, error) {
	t := &Text{
		Node: NewNode(node),
	}

	return t, t.Validate()
}

// QuickText creates a new Text from a string.
func QuickText(s string) *Text {
	n := htmx.TextNode(s)
	nn := NewNode(n)
	return &Text{
		Node: nn,
	}
}

// NewText creates a new Text from an html.Node.
func (p *Parser) NewText(node *html.Node) (*Text, error) {
	return NewText(node)
}
