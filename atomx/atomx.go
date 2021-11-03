package atomx

const (
	ERROR   Atom = "error"
	UNKNOWN Atom = "unknown"
)

type Atomable interface {
	Atom() Atom
}

type Atom string
type Atoms []Atom

func (a Atom) String() string {
	return string(a)
}

func (a Atom) Atom() Atom {
	return a
}

func (a Atom) Is(wants ...Atom) bool {
	for _, want := range wants {
		if a == want {
			return true
		}
	}
	return false
}

func IsAtom(a Atomable, wants ...Atom) bool {
	if a == nil {
		return false
	}

	at := a.Atom()
	return at.Is(wants...)
}
