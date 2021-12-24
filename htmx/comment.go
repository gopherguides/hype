package htmx

import (
	"golang.org/x/net/html"
)

// CommentNode returns a new html.CommentNode with the given text.
func CommentNode(text string) *html.Node {
	return &html.Node{
		Type: html.CommentNode,
		Data: text,
	}
}
