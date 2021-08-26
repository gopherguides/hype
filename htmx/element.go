package htmx

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func ElementNode(name string) *html.Node {
	return &html.Node{
		Data:     name,
		DataAtom: atom.Lookup([]byte(name)),
		Type:     html.ElementNode,
	}
}
