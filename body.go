package hype

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var _ Tag = &Body{}

type Body struct {
	*Node
}

func (b Body) String() string {
	sb := &strings.Builder{}
	sb.WriteString(b.StartTag())

	kids := b.GetChildren()
	if len(kids) > 0 {
		fmt.Fprintf(sb, "\n%s\n", kids)
	}

	sb.WriteString(b.EndTag() + "\n")
	return sb.String()
}

func (b Body) AsPage() *Page {
	p := &Page{
		Node: b.Clone(),
	}

	p.Data = "page"
	p.DataAtom = atom.Atom(0)

	return p
}

func (p *Parser) NewBody(node *Node) (*Body, error) {
	if node == nil || node.Node == nil {
		return nil, fmt.Errorf("body node can not be nil")
	}

	if node.Type != html.ElementNode {
		return nil, fmt.Errorf("node is not an element node %v", node)
	}

	b := &Body{
		Node: node,
	}

	return b, nil
}
