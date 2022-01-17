package hype

import (
	"encoding/json"
	"fmt"
	"io/fs"

	"github.com/gopherguides/hype/htmx"
	"golang.org/x/net/html"
)

// Document represents an HTML document
type Document struct {
	*Node
	fs.FS
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
	m := map[string]any{
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

	return ByType(doc.Children, &Meta{})
}

// Validate the document
func (d Document) Validate(p *Parser, checks ...ValidatorFn) error {
	chocks := ChildrenValidators(d, p, checks...)
	err := d.Node.Validate(p, html.DocumentNode, chocks...)

	return err
}

// NewDocument parses the node and returns a Document.
// The node must be of type html.DocumentNode.
func (p *Parser) NewDocument(n *html.Node) (*Document, error) {

	doc := &Document{
		FS:   p,
		Node: NewNode(n),
	}

	if err := doc.Validate(p); err != nil {
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

	err := doc.Validate(p)
	if err != nil {
		return nil, err
	}

	err = p.finalize(doc)
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

	// charset

	return nil
}

// Body returns the <body> tag for the document.
func (doc *Document) Body() (*Body, error) {
	if doc == nil {
		return nil, fmt.Errorf("document can not be nil")
	}

	bodies := ByType(doc.Children, &Body{})
	if len(bodies) == 0 {
		return nil, fmt.Errorf("body not found")
	}

	return bodies[0], nil
}

// Pages returns all of the <page> tags for the document.
func (doc *Document) Pages() Pages {
	if doc == nil {
		return nil
	}

	pages := ByType(doc.Children, &Page{})

	if len(pages) > 0 {
		return pages
	}

	body, err := doc.Body()
	if err != nil {
		return nil
	}
	return Pages{body.AsPage()}
}
