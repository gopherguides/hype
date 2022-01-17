package hype

import (
	"fmt"
	"strings"

	"github.com/gopherguides/hype/atomx"
	"golang.org/x/net/html"
)

var _ Tag = &Body{}
var _ Validatable = &Body{}

// Body represents the body of a document.
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

// AsPage returns the body as a Page.
func (b Body) AsPage() *Page {
	p := &Page{
		Node: b.Clone(),
	}

	p.DataAtom = atomx.Page

	return p
}

// Validate the body
func (b Body) Validate(p *Parser, checks ...ValidatorFn) error {
	checks = append(checks, AtomValidator("body"))
	return b.Node.Validate(p, html.ElementNode, checks...)
}

// NewBody returns a new Body from the given node.
func NewBody(node *Node) (*Body, error) {
	b := &Body{
		Node: node,
	}

	return b, nil
}
