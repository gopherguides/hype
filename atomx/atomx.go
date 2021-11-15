package atomx

import "strings"

const (
	ERROR   Atom = "error"
	UNKNOWN Atom = "unknown"
)

type Atomable interface {
	Atom() Atom
}

type Atom string

func (a Atom) String() string {
	return string(a)
}

func (a Atom) Atom() Atom {
	return a
}

func (a Atom) Is(wants ...Atom) bool {
	return Atoms(wants).Has(a)
}

func IsAtom(a Atomable, wants ...Atom) bool {
	if a == nil {
		return false
	}

	at := a.Atom()
	return Atoms(wants).Has(at)
}

type Atoms []Atom

func (atoms Atoms) String() string {
	ats := make([]string, 0, len(atoms))
	for _, at := range atoms {
		ats = append(ats, at.String())
	}
	return strings.Join(ats, ", ")
}

func (atoms Atoms) Has(a Atom) bool {
	for _, at := range atoms {
		if at == a {
			return true
		}
	}
	return false
}
