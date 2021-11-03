package hype

import (
	"fmt"
	"strings"

	"github.com/gopherguides/hype/atomx"
	"golang.org/x/net/html"
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

	p.DataAtom = atomx.Page

	return p
}

func (b Body) Validate(checks ...ValidatorFn) error {
	checks = append(checks, AdamValidator("body"))
	return b.Node.Validate(html.ElementNode, checks...)
}

func (p *Parser) NewBody(node *Node) (*Body, error) {
	return NewBody(node)
}

func NewBody(node *Node) (*Body, error) {
	b := &Body{
		Node: node,
	}

	return b, b.Validate()
}
