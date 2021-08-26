package htmx

import (
	"golang.org/x/net/html"
)

func CommentNode(text string) *html.Node {
	return &html.Node{
		Type: html.CommentNode,
		Data: text,
	}
}
