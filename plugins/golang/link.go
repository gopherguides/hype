package golang

import (
	"bytes"
	"fmt"
	"path"
	"strings"

	"github.com/gopherguides/hype"
	"github.com/gopherguides/hype/htmx"
	"golang.org/x/net/html"
)

var _ hype.Tag = &Link{}

type Link struct {
	*hype.Node
}

func (link Link) String() string {
	bb := &bytes.Buffer{}
	fmt.Fprint(bb, link.StartTag())
	fmt.Fprint(bb, link.Children.String())
	fmt.Fprint(bb, link.EndTag())
	return bb.String()
}

func (link *Link) Validate(checks ...hype.ValidatorFn) error {
	if link == nil {
		return fmt.Errorf("link is nil")
	}

	return link.Node.Validate(html.ElementNode, checks...)
}

func NewLink(node *hype.Node) (*Link, error) {
	link := &Link{
		Node: node,
	}

	if err := link.Validate(); err != nil {
		return nil, err
	}

	href, err := link.Get("href")

	if err != nil {
		if len(link.Children) == 0 {
			return nil, err
		}
		href = link.Children.String()
	}

	ats := link.Attrs()
	ats["href"] = "https://pkg.go.dev/" + href
	if _, ok := ats["target"]; !ok {
		ats["target"] = "_blank"
	}

	link.Set("for", href)

	dt := strings.ReplaceAll(href, "#", ".")

	if _, ok := ats["full-pkg"]; !ok {
		spl := strings.Split(href, "#")
		base := path.Base(spl[0])
		if len(spl) > 1 {
			dt = fmt.Sprintf("%s.%s", base, spl[1])
		}
	}

	a := &hype.Element{
		Node: hype.NewNode(htmx.AttrNode("a", ats)),
	}

	codeNode := hype.NewNode(htmx.AttrNode("code", nil))

	code, err := hype.NewInlineCode(codeNode)

	if err != nil {
		return nil, err
	}

	text := hype.QuickText(dt)
	code.Children = append(code.Children, text)
	a.Children = hype.Tags{code}

	link.Children = hype.Tags{a}

	return link, link.Validate()
}
