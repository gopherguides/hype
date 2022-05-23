package hype

// FindTitle finds the title element in the given HTML document.
// If no title element is found, the first h1 element is returned.
// If no h1 element is found, `"Untitled"` is returned.
func FindTitle(nodes Nodes) string {
	titles := ByAtom(nodes, "title")
	if len(titles) > 0 {
		return titles[0].Children().String()
	}

	h1s := ByAtom(nodes, "h1")
	if len(h1s) > 0 {
		return h1s[0].Children().String()
	}

	return "Untitled"
}
