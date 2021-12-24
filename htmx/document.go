package htmx

import (
	"golang.org/x/net/html"
)

// DocumetNode returns a new, empty html.DocumentNode.
func DocumentNode() *html.Node {
	return &html.Node{
		Type: html.DocumentNode,
	}
}
