package hype

import (
	"fmt"

	"golang.org/x/net/html"
)

var _ Tag = &Comment{}

// Comment represents an HTML comment.
type Comment struct {
	*Node
}

func (c Comment) String() string {
	return fmt.Sprintf("<!-- %s -->", c.Atom())
}

func (c Comment) Validate(checks ...ValidatorFn) error {
	return c.Node.Validate(html.CommentNode, checks...)
}

// NewComment returns a new Comment from the given node.
func (p *Parser) NewComment(node *html.Node) (*Comment, error) {
	return NewComment(node)
}

// NewComment returns a new Comment from the given node.
func NewComment(n *html.Node) (*Comment, error) {
	c := &Comment{
		Node: NewNode(n),
	}

	return c, c.Validate()
}
