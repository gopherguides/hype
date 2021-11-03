package hype

func findTitle(tags Tags) string {
	titles := tags.ByAdam("title")
	if len(titles) > 0 {
		return titles[0].GetChildren().String()
	}

	h1s := tags.ByAdam("h1")
	if len(h1s) > 0 {
		return h1s[0].GetChildren().String()
	}

	return "Untitled"
}
