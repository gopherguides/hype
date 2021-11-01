package hype

import "golang.org/x/net/html/atom"

func findTitle(tags Tags) string {
	titles := tags.ByAtom(atom.Title)
	if len(titles) > 0 {
		return titles[0].GetChildren().String()
	}

	h1s := tags.ByAtom(atom.H1)
	if len(h1s) > 0 {
		return h1s[0].GetChildren().String()
	}

	return "Untitled"
}
