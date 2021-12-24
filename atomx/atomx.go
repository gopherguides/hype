package atomx

import "strings"

const (
	ERROR   Atom = "error"
	UNKNOWN Atom = "unknown"
)

type Atomable interface {
	Atom() Atom
}

// Atom is a string that is used to identify a specific element or attribute.
// Example:
// 	"a", "img", "div"
type Atom string

// String representation of an Atom.
func (a Atom) String() string {
	return string(a)
}

// Atom implements the Atomable interface.
func (a Atom) Atom() Atom {
	return a
}

// Is returns true if the atom is in the list of atoms.
func (a Atom) Is(wants ...Atom) bool {
	return Atoms(wants).Has(a)
}

// IsAtom returns true if the atom is in the list of atoms.
func IsAtom(a Atomable, wants ...Atom) bool {
	if a == nil {
		return false
	}

	at := a.Atom()
	return Atoms(wants).Has(at)
}

// Atoms is a slice of Atoms.
type Atoms []Atom

// String returns a string representation of the atoms.
// Example:
//	"a, img, div"
func (atoms Atoms) String() string {
	ats := make([]string, 0, len(atoms))
	for _, at := range atoms {
		ats = append(ats, at.String())
	}
	return strings.Join(ats, ", ")
}

// Has returns true if the atom is in the list of atoms.
func (atoms Atoms) Has(a Atom) bool {
	for _, at := range atoms {
		if at == a {
			return true
		}
	}
	return false
}
