package hype

import (
	"fmt"
	"strings"

	"github.com/gopherguides/hype/atomx"
	"golang.org/x/net/html"
)

var _ Tag = &InlineCode{}
var _ Validatable = &InlineCode{}

// InlineCode represents inline code.
//
// Example:
// 	This is inline `code`.
// 	This is inline <code>code</code>.
type InlineCode struct {
	*Node
}

// Lang represents the language of an inline code snippet.
// Returns an empty string.
func (c *InlineCode) Lang() string {
	return ""
}

func (c *InlineCode) String() string {
	sb := &strings.Builder{}

	text := c.Children.String()
	text = strings.TrimSpace(text)

	fmt.Fprint(sb, c.StartTag())
	fmt.Fprint(sb, text)
	fmt.Fprint(sb, c.EndTag())
	return sb.String()
}

// Validate the InlineCode
func (inc InlineCode) Validate(p *Parser, checks ...ValidatorFn) error {
	checks = append(checks, AtomValidator(atomx.Code))
	return inc.Node.Validate(p, html.ElementNode, checks...)
}

// NewInlineCode returns a new InlineCode from the given node.
func NewInlineCode(node *Node) (*InlineCode, error) {
	c := &InlineCode{
		Node: node,
	}

	return c, c.Validate(nil)
}
