package htmx

import (
	"testing"

	"golang.org/x/net/html"
)

func TextNode(t *testing.T, text string) *html.Node {
	t.Helper()
	return &html.Node{
		Data: text,
		Type: html.TextNode,
	}
}
