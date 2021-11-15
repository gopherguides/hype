package hype

import (
	"fmt"
	"strings"

	"github.com/gopherguides/hype/atomx"
	"golang.org/x/net/html"
)

type InlineCode struct {
	*Node
}

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

func (inc InlineCode) Validate(checks ...ValidatorFn) error {
	checks = append(checks, AtomValidator(atomx.Code))
	return inc.Node.Validate(html.ElementNode, checks...)
}

func NewInlineCode(node *Node) (*InlineCode, error) {
	c := &InlineCode{
		Node: node,
	}

	return c, c.Validate()
}
