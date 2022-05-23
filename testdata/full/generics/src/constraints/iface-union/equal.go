package demo

import "golang.org/x/exp/constraints"

type Equalizer interface {
	constraints.Ordered
}

type Equalable interface {
	Equals(a any) bool
}

func Equal[T Equalizer](a T, b T) bool {
	return a == b
}
