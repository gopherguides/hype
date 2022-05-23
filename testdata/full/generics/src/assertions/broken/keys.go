package demo

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// snippet: def
func Keys[K constraints.Ordered, V any](m map[K]V) []K {
	// snippet: def

	// make a slice of the keys
	keys := make([]K, 0, len(m))

	// iterate over the map
	for k := range m {

		// if k implements fmt.Stringer,
		// print the string representation
		if st, ok := k.(fmt.Stringer); ok {
			fmt.Println(st.String())
		}

		// add the key to the slice
		keys = append(keys, k)
	}

	// return the keys
	return keys
}
