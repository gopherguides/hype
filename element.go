package hype

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type Element struct {
	*Node
}

func (e Element) String() string {
	switch e.Adam() {
	case "a", "img", "image", "link":
		return e.InlineTag()
	}

	sb := &strings.Builder{}
	sb.WriteString(e.StartTag())

	kids := e.GetChildren()
	if len(kids) > 0 {
		fmt.Fprintf(sb, "%s", kids)
	}

	sb.WriteString(e.EndTag())
	return sb.String()
}

func (e Element) Validate(checks ...ValidatorFn) error {
	return e.Node.Validate(html.ElementNode, checks...)
}

func (p *Parser) ElementNode(n *html.Node) (Tag, error) {
	node := NewNode(n)

	err := node.Validate(html.ElementNode)
	if err != nil {
		return nil, err
	}

	c := n.FirstChild
	for c != nil {
		tag, err := p.ParseNode(c)
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, tag)
		c = c.NextSibling
	}

	p.RLock()
	if ct := p.customTags; ct != nil {
		if fn, ok := ct[n.Data]; ok {
			p.RUnlock()
			return fn(node)
		}
	}
	p.RUnlock()

	switch node.Adam() {
	case "img", "image":
		return p.NewImage(node)
	case "meta":
		return p.NewMeta(node)
	case Code_Adam:
		return p.NewCode(node)
	case "body":
		return p.NewBody(node)
	case File_Adam:
		return p.NewFile(node)
	case File_Group_Adam:
		return p.NewFileGroup(node)
	case Include_Adam:
		return p.NewInclude(node)
	case Page_Adam:
		return p.NewPage(node)
	}

	el := &Element{
		Node: node,
	}

	return el, el.Validate()
}
