package htmltest

import (
	"testing"

	// . "github.com/gopherguides/hype"
	"golang.org/x/net/html"
)

func DocTypeNode(t *testing.T, value string) *html.Node {
	t.Helper()
	return &html.Node{
		Type: html.DoctypeNode,
		Data: value,
	}
}
