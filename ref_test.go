package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Ref_MarshalJSON(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	ref := &Ref{
		Element: NewEl("ref", nil),
	}
	ref.Nodes = append(ref.Nodes, Text("foo"))

	err := ref.Set("id", "foo")
	r.NoError(err)

	testJSON(t, "ref", ref)

}
