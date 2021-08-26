package htmx

import (
	"testing"

	// . "github.com/gopherguides/hype"
	"golang.org/x/net/html"
)

func DocumentNode(t *testing.T) *html.Node {
	t.Helper()

	return &html.Node{
		Type: html.DocumentNode,
	}
}
