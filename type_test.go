package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToType(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	r.Equal("hype.Text", toType(Text("x")))

	p := &Element{}
	r.Equal("hype.Element", toType(p))

	r.Equal("string", toType("hello"))
}
