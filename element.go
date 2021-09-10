package hype

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Element struct {
	*Node
}

func (e Element) String() string {
	switch e.DataAtom {
	case atom.Link, atom.Img:
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

func (p *Parser) ElementNode(node *html.Node) (Tag, error) {
	if node == nil {
		return nil, fmt.Errorf("node can not be nil")
	}

	if node.Type != html.ElementNode {
		return nil, fmt.Errorf("node is not an element node %v", node)
	}

	g := NewNode(node)
	c := node.FirstChild
	for c != nil {
		tag, err := p.ParseNode(c)
		if err != nil {
			return nil, err
		}
		g.Children = append(g.Children, tag)
		c = c.NextSibling
	}

	p.RLock()
	if ct := p.customTags; ct != nil {
		if fn, ok := ct[node.Data]; ok {
			p.RUnlock()
			return fn(g)
		}
	}
	p.RUnlock()

	switch node.DataAtom {
	case atom.Img, atom.Image:
		return p.NewImage(g)
	case atom.Meta:
		return p.NewMeta(g)
	case atom.Code:
		return p.NewCode(g)
	case atom.Body:
		return p.NewBody(g)
	default:
		switch node.Data {
		case "include":
			return p.NewInclude(g)
		case "page":
			return p.NewPage(g)
		}
	}

	return &Element{
		Node: g,
	}, nil
}
