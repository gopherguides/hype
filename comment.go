package hype

import (
	"fmt"

	"golang.org/x/net/html"
)

var _ Tag = &Comment{}

type Comment struct {
	*Node
}

func (c Comment) String() string {
	return fmt.Sprintf("<!-- %s -->", c.Data)
}

func (p *Parser) NewComment(node *html.Node) (*Comment, error) {
	return NewComment(node)
}

func NewComment(node *html.Node) (*Comment, error) {
	if node == nil {
		return nil, fmt.Errorf("node can not be nil")
	}

	if node.Type != html.CommentNode {
		return nil, fmt.Errorf("node is not a comment node %v", node)
	}

	return &Comment{
		Node: NewNode(node),
	}, nil
}
