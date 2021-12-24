package htmx

import (
	"encoding/json"

	"golang.org/x/net/html"
)

// NodeJSON is a JSON representation of an html.Node.
type NodeJSON struct {
	Atom       string     `json:"atom,omitempty"`
	Attributes Attributes `json:"attributes,omitempty"`
	Children   []NodeJSON `json:"children,omitempty"`
	Data       string     `json:"data,omitempty"`
	Namespace  string     `json:"namespace,omitempty"`
	Type       string     `json:"type,omitempty"`
}

// String returns the marshaled JSON representation of the NodeJSON.
func (node NodeJSON) String() string {
	b, _ := json.Marshal(node)
	return string(b)
}

// MarshalNode will marshal the given html.Node into JSON.
func MarshalNode(node *html.Node) ([]byte, error) {
	return json.Marshal(NewNodeJSON(node))
}

// NewNodeJSON returns a new NodeJSON for the given html.Node.
func NewNodeJSON(node *html.Node) NodeJSON {
	if node == nil {
		return NodeJSON{}
	}

	var ty string
	switch node.Type {
	case html.CommentNode:
		ty = "comment"
	case html.DoctypeNode:
		ty = "doctype"
	case html.DocumentNode:
		ty = "document"
	case html.TextNode:
		ty = "text"
	case html.ElementNode:
		ty = "element"
	case html.ErrorNode:
		ty = "error"
	}

	nj := NodeJSON{
		Atom:       node.DataAtom.String(),
		Attributes: NewAttributes(node),
		Data:       node.Data,
		Namespace:  node.Namespace,
		Type:       ty,
	}

	c := node.FirstChild
	for c != nil {
		nj.Children = append(nj.Children, NewNodeJSON(c))
		c = c.NextSibling
	}

	return nj
}
