package hype

import (
	"encoding/json"
	"fmt"

	"github.com/markbates/fsx"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
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
	tags := d.Children.AllData("overview")
	if len(tags) == 0 {
		return ""
	}
	return tags[0].GetChildren().String()
}

func (d Document) MarshalJSON() ([]byte, error) {
	m := jmap{
		"document": NewNodeJSON(d.Node.Node),
		"fs":       d.FS,
	}
	return json.Marshal(m)
}

// Meta returns all of the <meta> tags for the document.
func (doc *Document) Meta() Metas {
	if doc == nil {
		return nil
	}

	meta := doc.Children.AllType(&Meta{})
	res := make([]*Meta, 0, len(meta))
	for _, m := range meta {
		if md, ok := m.(*Meta); ok {
			res = append(res, md)
		}
	}

	return res
}

// NewDocument parses the node and returns a Document.
// The node must be of type html.DocumentNode.
func (p *Parser) NewDocument(node *html.Node) (*Document, error) {
	if node == nil {
		return nil, fmt.Errorf("node can not be nil")
	}

	if node.Type != html.DocumentNode {
		return nil, fmt.Errorf("node is not a document %v", node)
	}

	doc := &Document{
		FS:   p.FS,
		Node: NewNode(node),
	}

	c := doc.Node.FirstChild
	for c != nil {
		tag, err := p.ParseNode(c)
		if err != nil {
			return nil, err
		}
		doc.Children = append(doc.Children, tag)
		c = c.NextSibling
	}

	return doc, nil
}

func (doc *Document) Body() (*Body, error) {
	if doc == nil {
		return nil, fmt.Errorf("document can not be nil")
	}

	bodies := doc.Children.AllAtom(atom.Body)
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
	titles := doc.Children.AllAtom(atom.Title)
	if len(titles) > 0 {
		return titles[0].GetChildren().String()
	}

	h1s := doc.Children.AllAtom(atom.H1)
	if len(h1s) > 0 {
		return h1s[0].GetChildren().String()
	}

	return "Untitled"
}

// Pages returns all of the <page> tags for the document.
func (doc *Document) Pages() Pages {
	if doc == nil {
		return nil
	}

	pages := doc.Children.AllData("page")
	res := make(Pages, 0, len(pages))

	if len(pages) == 0 {
		body, err := doc.Body()
		if err != nil {
			return nil
		}
		return append(res, &Page{
			Node: body.Node,
		})
	}

	for _, m := range pages {
		if md, ok := m.(*Page); ok {
			res = append(res, md)
		}
	}

	return res
}
