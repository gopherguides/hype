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
	r.Error(err)

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
			if tt.err {
				r.Error(err)
				return
			}

			r.NoError(err)
			r.NotNil(m)
			r.Equal(tt.ats, m.Attrs())
		})
	}

}
