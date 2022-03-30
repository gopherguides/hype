package htmx

import "golang.org/x/net/html"

type NodeType html.NodeType

func (nt NodeType) String() string {
	switch html.NodeType(nt) {
	case html.ErrorNode:
		return "ErrorNode"
	case html.TextNode:
		return "TextNode"
	case html.DocumentNode:
		return "DocumentNode"
	case html.ElementNode:
		return "ElementNode"
	case html.CommentNode:
		return "CommentNode"
	case html.DoctypeNode:
		return "DoctypeNode"
	case html.RawNode:
		return "RawNode"
	}
	return "Unknown"
}

// const (
// )
