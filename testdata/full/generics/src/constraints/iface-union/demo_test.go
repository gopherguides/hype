package demo

import (
	"testing"
)

func Test_Equal(t *testing.T) {
	t.Parallel()

	es := "one"
	as := "one"

	if !Equal(es, as) {
		t.Fatalf("expected %v to equal %v", es, as)
	}

	ei := 1
	ai := 1

	if !Equal(ei, ai) {
		t.Fatalf("expected %v to equal %v", ei, ai)
	}

	ef := 1.2
	af := 1.2

	if !Equal(ef, af) {
		t.Fatalf("expected %v to equal %v", ef, af)
	}
}
