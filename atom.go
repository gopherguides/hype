package hype

import "github.com/gopherguides/hype/atomx"

type (
	Atom     = atomx.Atom
	Atomable = atomx.Atomable
	Atoms    = atomx.Atoms
)

func IsAtom(a Atomable, wants ...Atom) bool {
	return atomx.IsAtom(a, wants...)
}
