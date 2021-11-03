package hype

import (
	"golang.org/x/net/html/atom"
)

const (
	FileGroup_Atom atom.Atom = 452184562
	File_Atom      atom.Atom = 1421757657
	Include_Atom   atom.Atom = 1818455657
	Page_Atom      atom.Atom = 1818488942
)

type Atomable interface {
	Atom() atom.Atom
}

func IsAtom(a Atomable, want atom.Atom) bool {
	if a == nil {
		return false
	}
	return a.Atom() == want
}
