package atomx

import (
	"testing"
)

func Test_Atoms_String(t *testing.T) {
	t.Parallel()

	ats := Atoms{A, P, B}

	exp := `a, p, b`
	act := ats.String()

	if exp != act {
		t.Fatalf("expected %q, got %q", exp, act)
	}

}

func Test_Atoms_Has(t *testing.T) {
	t.Parallel()

	ats := Atoms{A, P, B}

	exp := true
	act := ats.Has(A)

	if exp != act {
		t.Fatalf("expected %v, got %v", exp, act)
	}

	exp = false
	act = ats.Has(I)

	if exp != act {
		t.Fatalf("expected %v, got %v", exp, act)
	}
}
