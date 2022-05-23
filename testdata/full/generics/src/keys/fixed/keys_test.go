package demo

import (
	"sort"
	"testing"
)

func Test_Keys(t *testing.T) {
	t.Parallel()

	// snippet: example
	// create a map with some values
	m := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}

	// create an interstitial map to pass to the function
	im := map[any]any{}

	// copy the map into the interstitial map
	for k, v := range m {
		im[k] = v
	}

	// get the keys
	keys := Keys(im)

	// create a slice to hold the keys as
	// integers for comparison
	act := make([]int, 0, len(keys))

	// copy the keys into the integer slice
	for _, k := range keys {
		// assert that the key is an int
		i, ok := k.(int)
		if !ok {
			t.Fatalf("expected type int, got %T", k)
		}

		act = append(act, i)
	}

	// sort the returned keys for comparison
	sort.Slice(act, func(i, j int) bool {
		return act[i] < act[j]
	})
	// snippet: example

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
