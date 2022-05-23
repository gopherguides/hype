package demo

import (
	"sort"
	"testing"
)

// snippet: example
func Test_Keys(t *testing.T) {
	t.Parallel()

	// create a map with some values
	m := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}

	// get the keys
	act := Keys(m)

	// sort the returned keys for comparison
	sort.Slice(act, func(i, j int) bool {
		return act[i] < act[j]
	})

	// set the expected values
	exp := []int{1, 2, 3}

	// assert the length of the actual and expected values
	al := len(act)
	el := len(exp)
	if al != el {
		t.Fatalf("expected %d, but got %d", el, al)
	}

	// loop through the expected values and
	// assert they are in the actual values
	for i, v := range exp {
		if v != act[i] {
			t.Fatalf("expected %d, but got %d", v, act[i])
		}
	}

}

// snippet: example
