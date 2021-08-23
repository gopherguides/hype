package htmltest

import (
	"testing"

	"golang.org/x/net/html"
)

func CommentNode(t *testing.T, text string) *html.Node {
	t.Helper()
	return &html.Node{
		Type: html.CommentNode,
		Data: text,
	}
}
