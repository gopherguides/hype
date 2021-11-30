package hype

import (
	"encoding/json"
	"fmt"

	"github.com/gopherguides/hype/htmx"
	"github.com/markbates/fsx"
	"golang.org/x/net/html"
)

// Document represents an HTML document
type Document struct {
	*Node
	FS *fsx.FS
}

func (d Document) String() string {
	return d.Children.String()
}

// Overview returns the contents of the first <overview> tag in the document.
func (d *Document) Overview() string {
	tags := d.Children.ByAtom("overview")
	if len(tags) == 0 {
		return ""
	}
	return tags[0].GetChildren().String()
}

func (d Document) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"document": htmx.NewNodeJSON(d.Node.html),
		"fs":       d.FS,
	}
	return json.Marshal(m)
}

// Meta returns all of the <meta> tags for the document.
func (doc *Document) Meta() Metas {
	if doc == nil {
		return nil
	}

	meta := doc.Children.ByType(&Meta{})
	res := make([]*Meta, 0, len(meta))
	for _, m := range meta {
		if md, ok := m.(*Meta); ok {
			res = append(res, md)
		}
	}

	return res
}

func (d Document) Validate(checks ...ValidatorFn) error {
	fn := func(n *Node) error {

		return nil
	}

	chocks := ChildrenValidators(d, checks...)
	chocks = append(chocks, fn)
	err := d.Node.Validate(html.DocumentNode, chocks...)

	return err
}

// NewDocument parses the node and returns a Document.
// The node must be of type html.DocumentNode.
func (p *Parser) NewDocument(n *html.Node) (*Document, error) {

	doc := &Document{
		FS:   p.FS,
		Node: NewNode(n),
	}

	if err := doc.Validate(); err != nil {
		return nil, err
	}

	c := doc.Node.html.FirstChild
	for c != nil {
		tag, err := p.ParseNode(c)
		if err != nil {
			return nil, err
		}
		doc.Children = append(doc.Children, tag)
		c = c.NextSibling
	}

	err := doc.Validate()
	if err != nil {
		return nil, err
	}

	err = p.finalize(doc.Children...)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (p *Parser) finalize(tags ...Tag) error {
	for _, tag := range tags {
		if f, ok := tag.(Finalizer); ok {
			if err := f.Finalize(p); err != nil {
				return err
			}
		}

		err := p.finalize(tag.GetChildren()...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (doc *Document) Body() (*Body, error) {
	if doc == nil {
		return nil, fmt.Errorf("document can not be nil")
	}

	bodies := doc.Children.ByAtom("body")
	if len(bodies) == 0 {
		return nil, fmt.Errorf("body not found")
	}

	body, ok := bodies[0].(*Body)
	if !ok {
		return nil, fmt.Errorf("node not a body %v", bodies[0])
	}

	return body, nil
}

// Title returns the <title> tag contents.
// If there is no <title> then the first <h1> is used.
// Default: Untitled
func (doc *Document) Title() string {
	return findTitle(doc.Children)
}

// Pages returns all of the <page> tags for the document.
func (doc *Document) Pages() Pages {
	if doc == nil {
		return nil
	}

	pages := doc.Children.ByAtom("page")
	res := make(Pages, 0, len(pages))

	if len(pages) == 0 {
		body, err := doc.Body()
		if err != nil {
			return nil
		}
		return append(res, body.AsPage())
	}

	for _, m := range pages {
		if md, ok := m.(*Page); ok {
			res = append(res, md)
		}
	}

	return res
}
