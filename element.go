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

	p.RLock()
	if ct := p.customTags; ct != nil {
		if fn, ok := ct[n.Data]; ok {
			p.RUnlock()
			return fn(node)
		}
	}
	p.RUnlock()

	switch n.DataAtom {
	case atom.Img, atom.Image:
		return p.NewImage(node)
	case atom.Meta:
		return p.NewMeta(node)
	case atom.Code:
		return p.NewCode(node)
	case atom.Body:
		return p.NewBody(node)
	default:
		switch n.Data {
		case "file":
			return p.NewFile(node)
		case "filegroup":
			return p.NewFileGroup(node)
		case "include":
			return p.NewInclude(node)
		case "page":
			return p.NewPage(node)
		}
	}

	el := &Element{
		Node: node,
	}

	return el, el.Validate()
}
