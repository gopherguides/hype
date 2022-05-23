package demo

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

type MapKey interface {
	constraints.Ordered | fmt.Stringer
}

// snippet: def
func Keys[K MapKey, V any](m map[K]V) []K {
	// snippet: def

	// make a slice of the keys
	keys := make([]K, 0, len(m))

	// iterate over the map
	for k := range m {

		// add the key to the slice
		keys = append(keys, k)
	}

	// return the keys
	return keys
}
