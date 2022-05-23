package demo

import (
	"fmt"
	"testing"
)

// Mapify returns a map of the given key
// and value slices.
func Mapify(k []any, v []any) (map[any]any, error) {
	if len(k) != len(v) {
		return nil, fmt.Errorf("key and value lengths do not match")
	}

	m := map[any]any{}

	for i := range k {
		m[k[i]] = v[i]
	}

	return m, nil
}

func Test_Echo(t *testing.T) {
	t.Parallel()

	k := []int{1, 2, 3}
	v := []string{"a", "b", "c"}

	m, err := Mapify(k, v)
	if err != nil {
		t.Fatal(err)
	}

	if len(m) != len(k) {
		t.Fatalf("expected %d elements, got %d", len(k), len(m))
	}

}
