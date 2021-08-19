package hype

import "golang.org/x/net/html/atom"

type Atomable interface {
	Atom() atom.Atom
}

func IsAtom(a Atomable, want atom.Atom) bool {
	if a == nil {
		return false
	}
	return a.Atom() == want
}
