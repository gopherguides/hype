package hype

import (
	"github.com/gopherguides/hype/atomx"
)

type CustomTagFn func(node *Node) (Tag, error)

type TagMap map[atomx.Atom]CustomTagFn
