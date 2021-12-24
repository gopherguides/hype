package hype

import (
	"github.com/gopherguides/hype/htmx"
)

type Attributes = htmx.Attributes

var NewAttributes = htmx.NewAttributes

// SrcAttr returns the source from the given attributes,
// if it exists.
//
// HTML Attributes:
// 	src: The source of the tag.
// 	skip-src: returns false even if the src exists
func SrcAttr(at Attributes) (Source, bool) {
	if _, ok := at["skip-src"]; ok {
		return "", false
	}

	s, ok := at["src"]
	return Source(s), ok
}

// TagSource returns the source of the tag, if it exists.
func TagSource(tag Tag) (Source, bool) {
	if sc, ok := tag.(Sourceable); ok {
		return sc.Source()
	}

	return SrcAttr(tag.Attrs())
}
