package demo

import (
	"fmt"
	"sort"
	"testing"
)

func Test_Keys(t *testing.T) {
	t.Parallel()

	// create a map with some values
	m := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}

	// snippet: example
	// create a function variable pointing
	// to the Keys function
	fn := Keys

	// get the keys
	act := fn(m)
	// snippet: example

	// sort the returned keys for comparison
	sort.Slice(act, func(i, j int) bool {
		return act[i] < act[j]
	})

	// set the expected values
	exp := []int{1, 2, 3}

	// assert the length of the actual and expected values
	if len(exp) != len(act) {
		t.Fatalf("expected len(%d), but got len(%d)", len(exp), len(act))
	}

	// assert the types of the actual and expected values
	at := fmt.Sprintf("%T", act)
	et := fmt.Sprintf("%T", exp)

	if at != et {
		t.Fatalf("expected type %s, but got type %s", et, at)
	}

	// loop through the expected values and
	// assert they are in the actual values
	for i, v := range exp {
		if v != act[i] {
			t.Fatalf("expected %d, but got %d", v, act[i])
		}
	}

}
