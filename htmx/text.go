package htmx

import (
	"golang.org/x/net/html"
)

// TextNode returns a new html.TextNode with the given text.
func TextNode(text string) *html.Node {
	return &html.Node{
		Data: text,
		Type: html.TextNode,
	}
}
