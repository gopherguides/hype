package hype

import (
	"fmt"

	"github.com/gopherguides/hype/atomx"
)

// Tag represents a tag in the HTML document.
type Tag interface {
	atomx.Atomable
	Attrs() Attributes
	GetChildren() Tags
	Nodeable
	fmt.Stringer
}

type StartTagger interface {
	StartTag() string
}

type EndTagger interface {
	EndTag() string
}

type Tagger interface {
	StartTagger
	EndTagger
}
