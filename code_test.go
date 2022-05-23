package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewCodeNodes_InlineCode(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	el := NewEl("code", nil)

	nodes, err := NewCodeNodes(nil, el)
	r.NoError(err)

	r.Len(nodes, 1)

	ic, ok := nodes[0].(*InlineCode)
	r.True(ok)
	r.Equal(ic.Element, el)
}

func Test_NewCodeNodes_SourceCode(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	el := NewEl("code", nil)

	r.NoError(el.Set("src", "main.go"))

	nodes, err := NewCodeNodes(nil, el)
	r.NoError(err)

	r.Len(nodes, 1)

	pre, ok := nodes[0].(*Element)
	r.True(ok)

	nodes = pre.Nodes
	r.Len(nodes, 1)

	sc, ok := nodes[0].(*SourceCode)
	r.True(ok)
	r.Equal(sc.Element, el)
}

func Test_NewCode_FencedCode(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	el := NewEl("code", nil)

	r.NoError(el.Set("language", "go"))

	nodes, err := NewCodeNodes(nil, el)
	r.NoError(err)

	r.Len(nodes, 1)

	fc, ok := nodes[0].(*FencedCode)
	r.True(ok)
	r.Equal(fc.Element, el)
}
