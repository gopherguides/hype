package hype

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	Page_Atom atom.Atom = 1818488942
)

var (
	_     Tag = &Page{}
	BREAK     = QuickText("<!--BREAK-->")
)

type Pages []*Page

type Page struct {
	*Node
}

func (c *Page) Src() string {
	c.RLock()
	defer c.RUnlock()
	return c.attrs["src"]
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

	fmt.Fprintln(sb, p.EndTag())
	return sb.String()
}

func (p *Page) EndTag() string {
	return fmt.Sprintf("%s%s", p.Node.EndTag(), BREAK)
}

func (p *Parser) NewPage(node *Node) (*Page, error) {
	return NewPage(node)
}

func NewPage(node *Node) (*Page, error) {
	if node == nil || node.Node == nil {
		return nil, fmt.Errorf("page node can not be nil")
	}

	if node.Type != html.ElementNode {
		return nil, fmt.Errorf("node is not an element node %v", node)
	}

	node.DataAtom = Page_Atom
	p := &Page{
		Node: node,
	}

	return p, nil
}
