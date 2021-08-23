package hype

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

var _ Tag = &Page{}

type Pages []*Page

type Page struct {
	*Node
}

func (p Page) Number() int {
	sn, _ := p.Get("number")
	i, _ := strconv.Atoi(sn)
	return i
}

func (p Page) Title() string {
	return findTitle(p.Children)
}

func (p Page) String() string {
	sb := &strings.Builder{}

	sb.WriteString(p.StartTag())

	kids := p.GetChildren()
	if len(kids) > 0 {
		fmt.Fprintf(sb, "\n%s\n", kids)
	}

	sb.WriteString(p.EndTag() + "\n")
	return sb.String()
}

func (p *Parser) NewPage(node *Node) (*Page, error) {
	if node == nil || node.Node == nil {
		return nil, fmt.Errorf("page node can not be nil")
	}

	if node.Type != html.ElementNode {
		return nil, fmt.Errorf("node is not an element node %v", node)
	}

	b := &Page{
		Node: node,
	}

	return b, nil
}
