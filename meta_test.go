package hype

import (
	"testing"

	"github.com/gopherguides/hype/htmx"
	"github.com/stretchr/testify/require"
)

func Test_NewMeta(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	_, err := NewMeta(nil)
	r.NoError(err)

	const key = `props`
	const val = `yo!`

	table := []struct {
		ats  Attributes
		err  bool
		name string
	}{
		{name: "name key", ats: Attributes{"name": key, "content": val}},
		{name: "property key", ats: Attributes{"property": key, "content": val}},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			n := htmx.AttrNode("meta", tt.ats)

			node := NewNode(n)

			m, err := NewMeta(node)

			r.NoError(err)
			r.NotNil(m)
			r.Equal(tt.ats, m.Attrs())
		})
	}

}
