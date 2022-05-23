package demo

import (
	"fmt"
)

// snippet: def
func Keys(m map[string]int) []string {
	// snippet: def

	// make a slice of the keys
	keys := make([]string, 0, len(m))

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
