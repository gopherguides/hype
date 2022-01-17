package hype

import (
	"bytes"
	"fmt"

	"github.com/gopherguides/hype/htmx"
	"golang.org/x/net/html"
)

var _ Tag = &Head{}
var _ Validatable = &Head{}

type Head struct {
	*Node
}

func (head Head) String() string {
	bb := &bytes.Buffer{}

	fmt.Fprint(bb, head.StartTag())
	fmt.Fprint(bb, head.GetChildren())
	fmt.Fprint(bb, head.EndTag())

	return bb.String()
}

func (head *Head) Validate(p *Parser, checks ...ValidatorFn) error {
	return head.Node.Validate(p, html.ElementNode, checks...)
}

func NewHead(node *Node) (*Head, error) {
	head := &Head{
		Node: node,
	}

	// metaFn := func() *Meta {
	// 	m := &Meta{
	// 		Node: NewNode(
	// 			htmx.AttrNode("meta", Attributes{
	// 				"charset": "utf-8",
	// 			}),
	// 		),
	// 	}
	// 	return m
	// }

	for _, m := range ByType(head.Children, &Meta{}) {
		if _, err := m.Get("charset"); err == nil {
			return head, nil
		}
	}

	m := &Meta{
		Node: NewNode(
			htmx.AttrNode("meta", Attributes{
				"charset": "utf-8",
			}),
		),
	}
	head.Children = append(head.Children, m)
	return head, nil
}

// heads := ByType(doc.Children.ByAtom(atomx.Head), &Element{})
// if len(heads) == 0 {
// 	return nil, fmt.Errorf("no <head> tag found")
// }

// head := heads[0]
// ms := head.Children.ByAttrs(Attributes{"charset": "*"})
// metas := ByType(ms, &Meta{})
// if len(metas) == 0 {
// 	head.Children = append(head.Children, &Meta{
// 		Node: NewNode(htmx.AttrNode("meta", Attributes{"charset": "utf-8"})),
// 	})
// }
