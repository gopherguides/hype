package atomx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Atoms_String(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	ats := Atoms{A, P, B}

	exp := `a, p, b`
	r.Equal(exp, ats.String())

}

func Test_Atoms_Has(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	ats := Atoms{A, P, B}
	r.True(ats.Has(A))
	r.False(ats.Has(I))
}
