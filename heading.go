package hype

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/gobuffalo/flect"
	"github.com/gopherguides/hype/atomx"
	"golang.org/x/net/html"
)

// Heading is an HTML heading element.
// 	H1, H2, H3, H4, H5, H6
//
// HTML Attributes:
// 	id: ID of the heading. Defaults to the heading text dasherized.
type Heading struct {
	*Node
	Parent *Heading // Parent heading
	Level  int      // Heading level
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
func (h Heading) Validate(checks ...ValidatorFn) error {
	checks = append(checks, AtomValidator(atomx.Headings()...))
	if len(h.ID()) == 0 {
		return fmt.Errorf("%s: missing id", h.Atom())
	}
	return h.Node.Validate(html.ElementNode, checks...)
}

// NewHeading returns a new Heading from the given node.
func NewHeading(node *Node) (*Heading, error) {
	h := &Heading{
		Node: node,
	}

	err := h.Validate()
	if err != nil {
		return nil, err
	}

	heads := atomx.Headings()
	for _, a := range heads {
		if a != node.Atom() {
			continue
		}

		s := strings.TrimPrefix(a.String(), "h")
		lvl, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		h.Level = lvl

	}

	return h, h.Validate()
}