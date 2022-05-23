package hype

import "github.com/gopherguides/hype/atomx"

type Atom = atomx.Atom

type Atomable interface {
	Atom() Atom
}

type AtomableNode interface {
	Node
	Atomable
}
