package htmx

import (
	"golang.org/x/net/html"
)

func DocTypeNode(value string) *html.Node {
	return &html.Node{
		Type: html.DoctypeNode,
		Data: value,
	}
}
