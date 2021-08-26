package htmx

import (
	"golang.org/x/net/html"
)

func TextNode(text string) *html.Node {
	return &html.Node{
		Data: text,
		Type: html.TextNode,
	}
}
