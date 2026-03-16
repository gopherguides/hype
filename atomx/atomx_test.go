package atomx

import (
	"testing"

	"github.com/stretchr/testify/require"
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

func Test_Atom_String(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	r.Equal("div", Div.String())
	r.Equal("a", A.String())
}

func Test_Atom_Atom(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	r.Equal(Div, Div.Atom())
	r.Equal(A, A.Atom())
}

func Test_Atom_Is(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	r.True(A.Is(A, B, P))
	r.False(A.Is(B, P, Div))
	r.False(A.Is())
}

func Test_IsAtom(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	r.False(IsAtom(nil, A))
	r.True(IsAtom(A, A, B))
	r.False(IsAtom(A, B, P))
}

func Test_Headings(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	h := Headings()
	r.Len(h, 6)
	r.True(h.Has(H1))
	r.True(h.Has(H2))
	r.True(h.Has(H3))
	r.True(h.Has(H4))
	r.True(h.Has(H5))
	r.True(h.Has(H6))
	r.False(h.Has(P))
}

func Test_Inlines(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	inl := Inlines()
	r.True(inl.Has(A))
	r.True(inl.Has(B))
	r.True(inl.Has(Br))
	r.True(inl.Has(Image))
	r.True(inl.Has(Img))
	r.True(inl.Has(Link))
	r.True(inl.Has(Ref))
	r.False(inl.Has(Div))
}
