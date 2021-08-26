package hype

import (
	"fmt"
	"strings"
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

func (p *Parser) NewInlineCode(node *Node) (*InlineCode, error) {
	return NewInlineCode(node)
}

func NewInlineCode(node *Node) (*InlineCode, error) {
	if node == nil || node.Node == nil {
		return nil, fmt.Errorf("inline code node can not be nil")
	}

	if node.Data != "code" {
		return nil, fmt.Errorf("node is not code %v", node.Data)
	}

	c := &InlineCode{
		Node: node,
	}

	return c, nil
}
