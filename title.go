package hype

import "fmt"

func findTitle(tags Tags) string {
	titles := tags.ByAtom("title")
	if len(titles) > 0 {
		return titles[0].GetChildren().String()
	}

	h1s := tags.ByAtom("h1")
	if len(h1s) > 0 {
		return h1s[0].GetChildren().String()
	}

	return "Untitled"
}

// Title returns the <title> tag contents.
// If there is no <title> then the first <h1> is used.
// Default: Untitled
func (doc *Document) Title() string {
	return findTitle(doc.Children)
}

func (doc *Document) SetTitle(title string) error {
	titles := doc.Children.ByAtom("title")
	if len(titles) > 0 {
		el, ok := titles[0].(*Element)
		if !ok {
			return fmt.Errorf("title is not an element, %T", titles[0])
		}

		el.Children = []Tag{QuickText(title)}
		return nil
	}

	h1s := doc.Children.ByAtom("h1")
	if len(h1s) > 0 {
		el, ok := h1s[0].(*Heading)
		if !ok {
			return fmt.Errorf("h1 is not a heading, %T", h1s[0])
		}
		el.Children = []Tag{QuickText(title)}
		return nil
	}

	return fmt.Errorf("no appriopriate title tag found")
}
