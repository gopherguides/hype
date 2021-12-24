package htmx

import (
	"golang.org/x/net/html"
)

// DocTypeNode returns a new html.Node with the given text.
func DocTypeNode(value string) *html.Node {
	return &html.Node{
		Type: html.DoctypeNode,
		Data: value,
	}
}
