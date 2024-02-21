package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FencedCode_MarshalJSON(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	code := &FencedCode{
		Element: NewEl("code", nil),
	}
	code.Nodes = append(code.Nodes, Text("var x = 1"))

	r.NoError(code.Set("language", "go"))

	testJSON(t, "fenced_code", code)

}
