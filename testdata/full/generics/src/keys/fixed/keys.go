package demo

// snippet: def
func Keys(m map[any]any) []any {
	// snippet: def

	// make a slice of the keys
	keys := make([]any, 0, len(m))

	// iterate over the map
	for k := range m {

		// add the key to the slice
		keys = append(keys, k)
	}

	// return the keys
	return keys
}