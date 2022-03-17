package hype

import (
	"bytes"
	"fmt"

	"github.com/gobuffalo/flect"
	"github.com/gopherguides/hype/atomx"
	"golang.org/x/net/html"
)

var _ Tag = &Heading{}
var _ Validatable = &Heading{}

// Heading is an HTML heading element.
// 	H1, H2, H3, H4, H5, H6
//
// HTML Attributes:
// 	id: ID of the heading. Defaults to the heading text dasherized.
type Heading struct {
	*Node
	Parent *Heading // Parent heading
}

func (h Heading) Title() string {
	return h.GetChildren().String()
}

func (h Heading) Level() int {
	for i, a := range atomx.Headings() {
		if h.DataAtom == a {
			return i + 1
		}
	}
	return 0
}

func (h Heading) String() string {
	bb := &bytes.Buffer{}
	fmt.Fprint(bb, h.StartTag())
	fmt.Fprint(bb, h.GetChildren())
	fmt.Fprint(bb, h.EndTag())
	return bb.String()
}

// ID returns the id of the heading.
func (h *Heading) ID() string {
	id, err := h.Get("id")
	if err == nil {
		return id
	}

	id = flect.Dasherize(h.GetChildren().String())
	if h.Parent != nil {
		id = h.Parent.ID() + "-" + id
	}
	return id
}

// Validate the heading
func (h Heading) Validate(p *Parser, checks ...ValidatorFn) error {
	checks = append(checks, AtomValidator(atomx.Headings()...))
	if len(h.ID()) == 0 {
		return fmt.Errorf("%s: missing id", h.Atom())
	}
	return h.Node.Validate(p, html.ElementNode, checks...)
}

// NewHeading returns a new Heading from the given node.
func NewHeading(node *Node) (*Heading, error) {
	h := &Heading{
		Node: node,
	}

	return h, nil
}
