package hype

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

var (
	_     Tag = &Page{}
	BREAK     = QuickText("<!--BREAK-->")
)

type Pages []*Page

// Page represents a page of a document.
type Page struct {
	*Node
}

// Title returns the title of the page.
// The title will be pulled from the first <title>
// tag on the page. If that does not exist,
// the first <h1> tag will be used. If the first
// <h1> tag does not exist, "untitled" will be
// returned.
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

func (p Page) Validate(checks ...ValidatorFn) error {
	return p.Node.Validate(html.ElementNode, checks...)
}

// NewPage returns a new Page from the given node.
func NewPage(node *Node) (*Page, error) {
	p := &Page{
		Node: node,
	}

	err := p.Validate()

	if err != nil {
		return nil, err
	}

	return p, p.Validate()
}
