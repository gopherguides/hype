package htmx

import (
	"testing"
	// . "github.com/gopherguides/hype"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func ElementNode(t *testing.T, name string) *html.Node {
	t.Helper()

	return &html.Node{
		Data:     name,
		DataAtom: atom.Lookup([]byte(name)),
		Type:     html.ElementNode,
	}
}
