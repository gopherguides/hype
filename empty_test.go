package hype

import (
	"testing"

	"github.com/gopherguides/hype/atomx"
	"github.com/stretchr/testify/require"
)

func Test_IsEmptyNode(t *testing.T) {
	t.Parallel()
	// r := require.New(t)

	// 	p := testParser(t, "testdata")

	// 	const frag = `<table>
	// <thead>
	// <tr>
	// <th></th>
	// <th></th>
	// </tr>
	// </thead>`

	// 	table, err := p.ParseFragment(strings.NewReader(frag))
	// 	r.NoError(err)

	fullP := NewEl(atomx.P, nil)
	fullP.Nodes = append(fullP.Nodes, Text("Hello World"))

	emptyP := NewEl(atomx.P, nil)

	tcs := []struct {
		name string
		node Node
		exp  bool
	}{
		{name: "nil", node: nil, exp: true},
		{name: "text", node: Text("hello"), exp: false},
		{name: "empty text", node: Text(""), exp: true},
		{name: "empty paragraph", node: emptyP, exp: true},
		{name: "full paragraph", node: fullP, exp: false},
		// {name: "empty table", node: table, exp: true},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := require.New(t)
			r.Equal(tc.exp, IsEmptyNode(tc.node))
		})
	}

}
