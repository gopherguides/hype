package hype

import (
	"encoding/json"

	"golang.org/x/net/html"
)

type NodeJSON struct {
	Atom       string     `json:"atom,omitempty"`
	Attributes Attributes `json:"attributes,omitempty"`
	Children   []NodeJSON `json:"children,omitempty"`
	Data       string     `json:"data,omitempty"`
	Namespace  string     `json:"namespace,omitempty"`
	Type       string     `json:"type,omitempty"`
}

func (node NodeJSON) String() string {
	b, _ := json.Marshal(node)
	return string(b)
}

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
