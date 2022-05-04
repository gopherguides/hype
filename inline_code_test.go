package hype

import (
	"testing"

	"github.com/gopherguides/hype/htmx"
	"github.com/stretchr/testify/require"
)

func Test_Parser_NewInlineCode(t *testing.T) {
	t.Parallel()

	valid := NewNode(htmx.ElementNode("code"))
	valid.Children = Tags{
		&Text{
			Node: NewNode(htmx.TextNode("hello")),
		},
	}

	table := []struct {
		err  bool
		exp  string
		name string
		node *Node
	}{
		{name: "nil", err: true},
		{name: "nil html node", node: &Node{}, err: true},
		{name: "non code node", node: NewNode(htmx.ElementNode("p")), err: true},
		{name: "valid", node: valid, exp: `<code class="inline-code">hello</code>`},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			sc, err := NewInlineCode(tt.node)
			if tt.err {
				r.Error(err)
				return
			}

			r.NoError(err)
			r.NotNil(sc)

			r.Equal(tt.exp, sc.String())
		})
	}

}
