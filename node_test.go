package hype

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Nodes_Delete(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	e1 := &Element{}
	e2 := &Element{}
	e3 := &Element{}

	nodes := Nodes{e1, e2, e3}

	nodes = nodes.Delete(e1)

	r.Equal(2, len(nodes))

	r.Equal(e2, nodes[0])
	r.Equal(e3, nodes[1])

}
