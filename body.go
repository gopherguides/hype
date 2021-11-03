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

func (b Body) Validate(checks ...ValidatorFn) error {
	checks = append(checks, AtomValidator(atom.Body))
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
