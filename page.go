package hype

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

var (
	_     Tag         = &Page{}
	_     Validatable = &Page{}
	BREAK             = QuickText("<!--BREAK-->")
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
func (page Page) Title() string {
	return findTitle(page.Children)
}

func (page Page) String() string {
	sb := &strings.Builder{}

	sb.WriteString(page.StartTag())

	kids := page.GetChildren()
	if len(kids) > 0 {
		fmt.Fprintf(sb, "\n%s\n", kids)
	}

	fmt.Fprintln(sb, page.EndTag())
	return sb.String()
}

func (page Page) EndTag() string {
	return fmt.Sprintf("%s%s", page.Node.EndTag(), BREAK)
}

func (page Page) Validate(p *Parser, checks ...ValidatorFn) error {
	return page.Node.Validate(p, html.ElementNode, checks...)
}

func (page *Page) ShiftHeadings(n int) {
	heads := ByType(page.Children, &Heading{})
	for _, h := range heads {
		lvl := h.Level()
		lvl += n
		h.DataAtom = Atom(fmt.Sprintf("h%d", lvl))
	}
}

// NewPage returns a new Page from the given node.
func NewPage(node *Node) (*Page, error) {
	page := &Page{
		Node: node,
	}

	return page, nil
}
