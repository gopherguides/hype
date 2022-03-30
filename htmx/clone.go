package htmx

import "golang.org/x/net/html"

// CloneNode returns a copy of the given node.
func CloneNode(n *html.Node) *html.Node {
	if n == nil {
		return nil
	}

	return &html.Node{
		Attr:        n.Attr,
		Data:        n.Data,
		DataAtom:    n.DataAtom,
		FirstChild:  CloneNode(n.FirstChild),
		LastChild:   CloneNode(n.LastChild),
		Namespace:   n.Namespace,
		NextSibling: CloneNode(n.NextSibling),
		Parent:      n.Parent,
		PrevSibling: n.PrevSibling,
		Type:        n.Type,
	}
}
