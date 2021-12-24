package hype

import (
	"github.com/gopherguides/hype/atomx"
)

// CustomTagFn is a function that returns a custom tag.
type CustomTagFn func(node *Node) (Tag, error)

// TagMap is a map of custom tags.
type TagMap map[atomx.Atom]CustomTagFn
