package hype

import (
	"fmt"
	"strings"

	"github.com/gopherguides/hype/atomx"
	"golang.org/x/net/html"
)

var _ Tag = &Element{}
var _ Validatable = &Element{}

// Element is a generic HTML element.
type Element struct {
	*Node
}

func (el Element) Markdown() string {
	switch el.Atom() {
	case atomx.A:
		return fmt.Sprintf("[%s](%s)", el.GetChildren().Markdown(), el.attrs["href"])
	case atomx.B, atomx.Strong:
		return fmt.Sprintf("**%s**", el.GetChildren().Markdown())
	}

	return el.GetChildren().Markdown()
}

func (e Element) String() string {
	at := e.Atom()

	if at.Is(atomx.Inlines()...) {
		return e.InlineTag()
	}

	sb := &strings.Builder{}
	sb.WriteString(e.StartTag())

	kids := e.GetChildren()
	if len(kids) > 0 {
		fmt.Fprintf(sb, "%s", kids)
	}

	sb.WriteString(e.EndTag())
	return sb.String()
}

// Validate the element
func (e Element) Validate(p *Parser, checks ...ValidatorFn) error {
	return e.Node.Validate(p, html.ElementNode, checks...)
}

// NewElement returns an element node from the given node.
func (p *Parser) NewElement(n *html.Node) (Tag, error) {
	node := NewNode(n)

	err := node.Validate(p, html.ElementNode)
	if err != nil {
		return nil, err
	}

	c := n.FirstChild
	for c != nil {
		tag, err := p.ParseNode(c)
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, tag)
		c = c.NextSibling
	}

	if fn, ok := p.CustomTag(atomx.Atom(n.Data)); ok {
		return fn(node)
	}

	el := &Element{
		Node: node,
	}

	return el, nil
}
