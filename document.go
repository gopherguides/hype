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

func (d *Document) String() string {
	return d.Children.String()
}

func (d *Document) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"document": d.Node,
		"fs":       d.FS,
	}
	return json.MarshalIndent(m, "", "  ")
}

func (doc *Document) MetaData() MetaData {
	if doc == nil {
		return MetaData{}
	}
	return doc.Children.MetaData()
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

	var html Tag
	for _, tag := range doc.Children {
		if IsAtom(tag, atom.Html) {
			html = tag
			break
		}
	}

	if html == nil {
		return nil, fmt.Errorf("no body found %v", doc)
	}

	for _, tag := range html.GetChildren() {
		if b, ok := tag.(*Body); ok {
			return b, nil
		}
	}

	return nil, fmt.Errorf("no body found %v", doc)
}
