package hype

import (
	"fmt"
	"strings"
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

// func (p *Parser) NewBody(node *html.Node) (*Body, error) {
// 	if node == nil {
// 		return nil, fmt.Errorf("node can not be nil")
// 	}

// 	if node.Type != html.ElementNode {
// 		return nil, fmt.Errorf("node is not an element node %v", node)
// 	}

// 	b := &Body{
// 		Node: NewNode(node),
// 	}

// 	return b, nil
// }
