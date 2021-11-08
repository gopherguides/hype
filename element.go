package hype

import (
	"fmt"
	"strings"

	"github.com/gopherguides/hype/atomx"
	"golang.org/x/net/html"
)

type Element struct {
	*Node
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

func (e Element) Validate(checks ...ValidatorFn) error {
	return e.Node.Validate(html.ElementNode, checks...)
}

func (p *Parser) ElementNode(n *html.Node) (Tag, error) {
	node := NewNode(n)

	err := node.Validate(html.ElementNode)
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

	fn, ok := p.CustomTag(atomx.Atom(n.Data))
	if ok {
		return fn(node)
	}

	switch node.Atom() {

	case atomx.File:
		return p.NewFile(node)
	case atomx.Filegroup:
		return p.NewFileGroup(node)
	case atomx.Include:
		return p.NewInclude(node)
	case atomx.Page:
		return p.NewPage(node)
	}

	el := &Element{
		Node: node,
	}

	return el, el.Validate()
}
