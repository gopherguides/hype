package htmx

import (
	"golang.org/x/net/html"
)

func DocumentNode() *html.Node {
	return &html.Node{
		Type: html.DocumentNode,
	}
}
