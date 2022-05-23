package demo

import "golang.org/x/exp/constraints"

type Model[T constraints.Ordered] interface {
	ID() T
}
