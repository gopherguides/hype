package hype

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	_     Tag = &Page{}
	BREAK     = QuickText("<!--BREAK-->")
)

type Pages []*Page

type Page struct {
	*Node
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

func (p *Page) Atom() atom.Atom {
	p.Lock()
	defer p.Unlock()

	if p.Node.DataAtom != Page_Atom {
		p.Node.DataAtom = Page_Atom
	}

	return Page_Atom
}

func (p *Page) EndTag() string {
	return fmt.Sprintf("%s%s", p.Node.EndTag(), BREAK)
}

func (p Page) Validate(checks ...ValidatorFn) error {
	return p.Node.Validate(html.ElementNode, checks...)
}

func NewPage(node *Node) (*Page, error) {
	p := &Page{
		Node: node,
	}

	err := p.Validate()

	if err != nil {
		return nil, err
	}

	node.DataAtom = Page_Atom

	return p, p.Validate()
}

func (p *Parser) NewPage(node *Node) (*Page, error) {
	return NewPage(node)
}
