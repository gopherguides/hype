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

	// if src, ok := node.attrs["src"]; ok {
	// 	fmt.Printf("TODO >> element.go:56 p.Root %[1]T %[1]v\n", p.Root)
	// 	fmt.Printf("TODO >> element.go:57 src %[1]T %[1]v\n", src)
	// 	src = filepath.Join(p.Root, src)
	// 	fmt.Printf("TODO >> element.go:59 src %[1]T %[1]v\n", src)
	// 	node.attrs["src"] = src
	// }

	if fn, ok := p.CustomTag(atomx.Atom(n.Data)); ok {
		return fn(node)
	}

	el := &Element{
		Node: node,
	}

	return el, el.Validate()
}
