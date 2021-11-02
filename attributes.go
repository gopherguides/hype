package hype

import (
	"github.com/gopherguides/hype/htmx"
)

type Attributes = htmx.Attributes

var NewAttributes = htmx.NewAttributes

func SrcAttr(at Attributes) (Source, bool) {
	if _, ok := at["skip-src"]; ok {
		return "", false
	}

	s, ok := at["src"]
	return Source(s), ok
}

func TagSource(tag Tag) (Source, bool) {
	if sc, ok := tag.(Sourceable); ok {
		return sc.Source()
	}

	return SrcAttr(tag.Attrs())
}
